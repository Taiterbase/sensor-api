package api

import (
	"context"
	"net/http"
	"time"
)

// Middleware used for CORS, authentication, logging, etc.

// TimeoutMiddleware adds a timeout to the request context.
func TimeoutMiddleware(timeout time.Duration, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		r = r.WithContext(ctx)
		done := make(chan struct{})
		go func() {
			h.ServeHTTP(w, r)
			close(done)
		}()

		select {
		case <-ctx.Done():
			w.WriteHeader(http.StatusRequestTimeout)
			w.Write([]byte("Request timeout"))
		case <-done:
		}
	})
}
