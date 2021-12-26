package middleware

import "net/http"

// Json middleware intercepts all http handlers and sets the content-type to
// be of json application. This is header will be included in the response.
// This tells the requester what format is the response so that they can format
// the response message accordingly.
func Json(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
