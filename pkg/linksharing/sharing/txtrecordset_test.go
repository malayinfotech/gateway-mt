// Copyright (C) 2020 Storx Labs, Inc.
// See LICENSE for copying information.

package sharing

import (
	"testing"
	"time"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/require"
)

func TestFull(t *testing.T) {
	testResponse := &dns.Msg{
		MsgHdr: dns.MsgHdr{
			Id: 0xa2c, Response: true, Opcode: 0, Authoritative: false, Truncated: false,
			RecursionDesired: true, RecursionAvailable: true, Zero: false,
			AuthenticatedData: false, CheckingDisabled: false, Rcode: 0},
		Compress: false,
		Question: []dns.Question{
			{
				Name:  "txt-staging.oliofam.test.",
				Qtype: 0x10, Qclass: 0x1}},
		Answer: []dns.RR{
			&dns.TXT{
				Hdr: dns.RR_Header{
					Name:   "txt-staging.oliofam.test.",
					Rrtype: 0x10, Class: 0x1, Ttl: 0x111, Rdlength: 0x1c},
				Txt: []string{
					"storx-root:oliosite/staging"}},
			&dns.TXT{
				Hdr: dns.RR_Header{Name: "txt-staging.oliofam.test.",
					Rrtype: 0x10, Class: 0x1, Ttl: 0x110, Rdlength: 0x45}, // 0x110 is 272
				Txt: []string{
					"storx-Access-1:1Z6ZMX4RM1oLAfJkcsXboSeuZ1A8wWUvYD8oQbMGhfK1FG17uS7Pj"}},
			&dns.TXT{
				Hdr: dns.RR_Header{Name: "txt-staging.oliofam.test.",
					Rrtype: 0x10, Class: 0x1, Ttl: 0x111, Rdlength: 0x45},
				Txt: []string{
					"storx_access-2:jCinS951sTrrpPFFEWH9zEzQCE5iCUR4iGMKRnGCsJbqFn8SMM4cd"}},
			&dns.TXT{
				Hdr: dns.RR_Header{Name: "txt-staging.oliofam.test.",
					Rrtype: 0x10, Class: 0x1, Ttl: 0x111, Rdlength: 0x45},
				Txt: []string{
					"storx-access-3:Y7XEntEDyKE3Q2woUoM7D6xHqD4RJFApee9BSpisLqVhC8YUs7Vvt"}},
			&dns.TXT{
				Hdr: dns.RR_Header{Name: "txt-staging.oliofam.test.",
					Rrtype: 0x10, Class: 0x1, Ttl: 0x111, Rdlength: 0x45},
				Txt: []string{
					"storx-access-4:NTiQYrvmvbzb4CLs2TryCxssFVuLrruDL88rQefHbJRUEnc34dRpi"}},
			&dns.TXT{
				Hdr: dns.RR_Header{Name: "txt-staging.oliofam.test.",
					Rrtype: 0x10, Class: 0x1, Ttl: 0x111, Rdlength: 0x92},
				Txt: []string{
					"storx-access-5:v9JuVwnpFpnyV275iaUdyEFAca8D7gCUr8ebbGQMHZzZVfzN86mt4",
					"storx-access-6:khSdyk1tCkAnP6sJ6reNZ5AfX89N7Yf4vv7fqXN5oPLL2iLqZ2Y8u",
					"jt-test:2"}},
			&dns.TXT{
				Hdr: dns.RR_Header{Name: "txt-staging.oliofam.test.",
					Rrtype: 0x10, Class: 0x1, Ttl: 0x111, Rdlength: 0x8a},
				Txt: []string{
					"storx-access-7:DDv46GqzsCd2Tfv3FjY7qXBwVUubGVUnyU5BMpdu21EgEo65jkpr8",
					"storx-access-8:JLfsx1zpLWbiJPsfQceVTJH2gF3h3eaVL1toHw9x2mSx4mumwuDac"}},
			&dns.TXT{
				Hdr: dns.RR_Header{Name: "txt-staging.oliofam.test.",
					Rrtype: 0x10, Class: 0x1, Ttl: 0x111, Rdlength: 0x1c},
				Txt: []string{"storx-access-9:4abuGDg4zL3j"}},
		},
		Ns:    []dns.RR(nil),
		Extra: []dns.RR(nil),
	}

	set := ResponseToTXTRecordSet(testResponse)
	require.Equal(t, set.Lookup("storx-grant"), "")
	require.Equal(t, set.Lookup("storx-access"), "1Z6ZMX4RM1oLAfJkcsXboSeuZ1A8wWUvYD8oQbMGhfK1FG17uS7PjjCinS951sTrrpPFFEWH9zEzQCE5iCUR4iGMKRnGCsJbqFn8SMM4cdY7XEntEDyKE3Q2woUoM7D6xHqD4RJFApee9BSpisLqVhC8YUs7VvtNTiQYrvmvbzb4CLs2TryCxssFVuLrruDL88rQefHbJRUEnc34dRpiv9JuVwnpFpnyV275iaUdyEFAca8D7gCUr8ebbGQMHZzZVfzN86mt4khSdyk1tCkAnP6sJ6reNZ5AfX89N7Yf4vv7fqXN5oPLL2iLqZ2Y8uDDv46GqzsCd2Tfv3FjY7qXBwVUubGVUnyU5BMpdu21EgEo65jkpr8JLfsx1zpLWbiJPsfQceVTJH2gF3h3eaVL1toHw9x2mSx4mumwuDac4abuGDg4zL3j")
	require.Equal(t, set.Lookup("storx-root"), "oliosite/staging")
	require.Equal(t, set.TTL(), 272*time.Second)
}
