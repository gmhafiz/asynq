package email

import "net/http"

type Health interface {
	Send(w http.ResponseWriter, r *http.Request)
}
