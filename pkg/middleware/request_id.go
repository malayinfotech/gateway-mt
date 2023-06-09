// Copyright (C) 2022 Storx Labs, Inc.
// See LICENSE for copying information.

package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/spacemonkeygo/monkit/v3"
)

// RequestIDKey is the key that holds the unique request ID in a request context.
type requestIDKey struct{}

// XStorxRequestID is the header key for the request ID.
const XStorxRequestID = "X-Storx-Request-Id"

// AddRequestID uses XStorxRequestID to set a unique request ID in response
// headers for each request if it doesn't already exist.
func AddRequestID(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(XStorxRequestID)
		if requestID == "" {
			requestID = fmt.Sprintf("%x", monkit.NewId())
		}

		w.Header().Set(XStorxRequestID, requestID)
		ctx := context.WithValue(r.Context(), requestIDKey{}, requestID)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRequestID returns the request ID from the context.
func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if reqID, ok := ctx.Value(requestIDKey{}).(string); ok {
		return reqID
	}
	return ""
}

// AddRequestIDToHeaders adds the request ID from the context to the request header.
func AddRequestIDToHeaders(req *http.Request) {
	if req == nil {
		return
	}

	// Ideally, the context should always have request ID, since it is being set in the middleware.
	req.Header.Set(XStorxRequestID, GetRequestID(req.Context()))
}
