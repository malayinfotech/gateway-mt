// Copyright (C) 2020 Storx Labs, Inc.
// See LICENSE for copying information.

package authdb

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"common/encryption"
	"common/grant"
	"common/macaroon"
	"common/storx"
	"common/testcontext"
)

func TestBase32(t *testing.T) {
	oldKey, err := NewEncryptionKey()
	require.NoError(t, err)
	require.NotEqual(t, oldKey, EncryptionKey{})
	encoded := oldKey.ToBase32()
	require.Len(t, encoded, EncKeySizeEncoded)
	var newKey EncryptionKey
	err = newKey.FromBase32(encoded)
	require.NoError(t, err)
	require.Equal(t, newKey, oldKey)
	require.NotEqual(t, newKey, EncryptionKey{})
}

func TestBase32Fail(t *testing.T) {
	key, err := NewEncryptionKey()
	require.NoError(t, err)
	encoded := key.ToBase32()
	require.Len(t, encoded, EncKeySizeEncoded)
	require.Error(t, key.FromBase32(encoded[1:]))
	require.Error(t, key.FromBase32(encoded[2:]))
	encoded = "a" + encoded[1:]
	require.Error(t, key.FromBase32(encoded))
}

func TestNonceIncrement(t *testing.T) {
	nonce := storx.Nonce{}
	_, err := encryption.Increment(&nonce, 1)
	require.NoError(t, err)
	require.Equal(t, storx.Nonce{1}, nonce)
}

func combineNotAfterCaveats(t *testing.T, unrestricted *macaroon.APIKey, times ...time.Time) *macaroon.APIKey {
	var (
		restricted = unrestricted
		err        error
	)

	for i, time := range times {
		restricted, err = restricted.Restrict(macaroon.Caveat{NotAfter: &time})
		require.NoError(t, err)
		if i%2 == 0 { // add noise
			restricted, err = restricted.Restrict(macaroon.Caveat{AllowedPaths: []*macaroon.Caveat_Path{{Bucket: []byte(strconv.Itoa(i))}}})
			require.NoError(t, err)
		}
	}

	return restricted
}

func TestApiKeyExpiration(t *testing.T) {
	unrestricted, err := macaroon.NewAPIKey([]byte("test"))
	require.NoError(t, err)

	withUnrelatedCaveats, err := unrestricted.Restrict(macaroon.Caveat{DisallowReads: true})
	require.NoError(t, err)
	withUnrelatedCaveats, err = withUnrelatedCaveats.Restrict(macaroon.Caveat{NotBefore: &time.Time{}})
	require.NoError(t, err)

	// a, b and c are times in the order of appearance:
	a := time.Now()
	b := a.Add(1)
	c := b.Add(1)

	tests := [...]struct {
		apiKey *macaroon.APIKey
		want   *time.Time
	}{
		{unrestricted, nil},
		{withUnrelatedCaveats, nil},
		{combineNotAfterCaveats(t, unrestricted, a), &a},
		{combineNotAfterCaveats(t, unrestricted, a, b), &a},
		{combineNotAfterCaveats(t, unrestricted, b, a), &a},
		{combineNotAfterCaveats(t, unrestricted, a, b, c), &a},
		{combineNotAfterCaveats(t, unrestricted, a, c, b), &a},
		{combineNotAfterCaveats(t, unrestricted, b, a, c), &a},
		{combineNotAfterCaveats(t, unrestricted, b, c, a), &a},
		{combineNotAfterCaveats(t, unrestricted, c, a, b), &a},
		{combineNotAfterCaveats(t, unrestricted, c, b, a), &a},
	}

	for i, tt := range tests {
		got, err := apiKeyExpiration(tt.apiKey)
		require.NoError(t, err, i)
		if tt.want != nil {
			require.Equal(t, tt.want.UTC(), got.UTC(), i)
		} else {
			require.Equal(t, tt.want, got, i)
		}
	}
}

func TestApiKeyExpiration_Invalid(t *testing.T) {
	mac, err := macaroon.NewUnrestricted([]byte("test"))
	require.NoError(t, err)
	mac, err = mac.AddFirstPartyCaveat([]byte("\xff\xfftrash\xff\xff"))
	require.NoError(t, err)

	k, err := macaroon.ParseRawAPIKey(mac.Serialize())
	require.NoError(t, err)

	_, err = apiKeyExpiration(k) // first caveat is invalid
	require.Error(t, err)

	k, err = k.Restrict(macaroon.Caveat{NotAfter: &time.Time{}})
	require.NoError(t, err)
	k, err = k.Restrict(macaroon.Caveat{DisallowDeletes: true})
	require.NoError(t, err)

	_, err = apiKeyExpiration(k) // one of the caveats is invalid
	require.Error(t, err)
}
func TestPutSatelliteValidation(t *testing.T) {
	ctx := testcontext.New(t)
	defer ctx.Cleanup()
	validURL := "12EayRS2V1kEsWESU9QMRseFhdxYxKicsiFmxrsLZHeLUtdps3S@us1.storx.io:7777"
	checkURL := "12EayRS2V1kEsWESU9QMRseFhdxYxKicsiFmxrsLZHeLUtdps3S@127.0.0.1:7778"

	mac, err := macaroon.NewAPIKey(nil)
	require.NoError(t, err)

	g := grant.Access{
		SatelliteAddress: checkURL,
		EncAccess:        grant.NewEncryptionAccess(),
		APIKey:           mac,
	}
	invalidGrant, err := g.Serialize()
	require.NoError(t, err)

	g.SatelliteAddress = validURL
	validGrant, err := g.Serialize()
	require.NoError(t, err)

	url, err := storx.ParseNodeURL(validURL)
	require.NoError(t, err)
	db := NewDatabase(mockKV{}, map[storx.NodeURL]struct{}{url: {}})

	key, err := NewEncryptionKey()
	require.NoError(t, err)

	_, err = db.Put(ctx, key, validGrant, false)
	require.NoError(t, err)
	_, err = db.Put(ctx, key, invalidGrant, false)
	require.Error(t, err)
}

type mockKV struct{}

func (mockKV) Put(ctx context.Context, keyHash KeyHash, record *Record) (err error) { return nil }
func (mockKV) Get(ctx context.Context, keyHash KeyHash) (record *Record, err error) { return nil, nil }
func (mockKV) PingDB(ctx context.Context) error                                     { return nil }
func (mockKV) Run(ctx context.Context) error                                        { return nil }
func (mockKV) Close() error                                                         { return nil }
