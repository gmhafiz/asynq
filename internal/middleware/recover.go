package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"

	chiMiddleware "github.com/go-chi/chi/middleware"
)

// Recovery adapted from https://github.com/go-chi/chi/blob/master/middleware/recoverer.go
// Panic can still happen. This middleware intercepts all panics and handles
// them by logging, and optionally send the error somewhere else.
//
// Handling them prevents the whole API from shutting down.
func Recovery(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {

				logEntry := chiMiddleware.GetLogEntry(r)
				if logEntry != nil {
					logEntry.Panic(rvr, debug.Stack())
				} else {
					chiMiddleware.PrintPrettyStack(rvr)
				}

				log.Println("PANIC")
				// notify to sentry/email/slack
				log.Printf("request: %s %s\n", r.Method, r.URL.RequestURI())
				body := r.GetBody
				b, _ := json.Marshal(body)
				log.Printf("host: %s\n", r.Host)
				log.Printf("request body: %s\n", b)

				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
