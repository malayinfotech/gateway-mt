// Copyright (C) 2020 Storx Labs, Inc.
// See LICENSE for copying information.

package auth

import (
	"context"

	"github.com/zeebo/errs"
	"go.uber.org/zap"

	"gateway-mt/pkg/auth/authdb"
	"gateway-mt/pkg/auth/badgerauth"
	"private/dbutil"
)

// OpenKV opens the database connection with the appropriate driver.
func OpenKV(ctx context.Context, log *zap.Logger, config Config) (_ authdb.KV, err error) {
	defer mon.Task()(&ctx)(&err)

	driver, _, _, err := dbutil.SplitConnStr(config.KVBackend)
	if err != nil {
		return nil, err
	}

	switch driver {
	case "badger":
		return badgerauth.New(log, config.Node)
	default:
		return nil, errs.New("unknown scheme: %q", config.KVBackend)
	}
}
