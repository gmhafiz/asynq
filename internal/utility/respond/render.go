package respond

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

var (
	ErrBadRequest    = errors.New("bad request")
	ErrNoRecord      = errors.New("no record found")
	ErrInternalError = errors.New("internal")

	ErrDatabase       = errors.New("connecting to database")
	ErrInvalidRequest = errors.New("invalid request")
)

func Render(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.WriteHeader(statusCode)

	if payload == nil {
		_, err := w.Write(nil)
		if err != nil {
			log.Println(err)
		}
	} else {
		data, err := json.Marshal(payload)
		if err != nil {
			log.Println(err)
		}
		_, err = w.Write(data)
		if err != nil {
			log.Println(err)
		}
	}
}

func Errors(w http.ResponseWriter, statusCode int, errors []string) {
	w.WriteHeader(statusCode)

	if errors == nil {
		write(w, nil)
		return
	}

	p := map[string][]string{
		"message": errors,
	}
	data, err := json.Marshal(p)
	if err != nil {
		log.Println(err)
	}

	if string(data) == "null" {
		return
	}

	write(w, data)
}

func Error(w http.ResponseWriter, statusCode int, message error) {
	w.WriteHeader(statusCode)

	var p map[string]string
	if message == nil {
		write(w, nil)
		return
	}

	p = map[string]string{
		"message": message.Error(),
	}
	data, err := json.Marshal(p)
	if err != nil {
		log.Println(err)
	}

	if string(data) == "null" {
		return
	}

	write(w, data)
}

func write(w http.ResponseWriter, data []byte) {
	_, err := w.Write(data)
	if err != nil {
		log.Println(err)
	}
}
