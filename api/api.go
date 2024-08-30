package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type SuccessResponse struct {
	Data []any `json:"data"`
}

func SuccessJson(w http.ResponseWriter, r *http.Request, data interface{}) {
	jsonMsg, err := json.Marshal(data)
	if err != nil {
		Error(w, r, fmt.Errorf("serialising response failed: %w", err), 500)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		success(w, r, jsonMsg)
	}
}

func success(w http.ResponseWriter, r *http.Request, data []byte) {
	if _, err := w.Write(data); err != nil {
		log.Printf("Error writing response: %v", err)
	}

	log.Printf(
		"%s %s %s 200",
		r.Method,
		r.RequestURI,
		r.RemoteAddr,
	)
}

func Error(w http.ResponseWriter, r *http.Request, err error, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err == nil {
		err = fmt.Errorf("nil err")
	}
	logErr := err

	errorMsgJSON, err := json.Marshal(ErrorResponse{Message: err.Error()})

	if err != nil {
		log.Println(err)
	} else {
		if _, err = w.Write(errorMsgJSON); err != nil {
			log.Printf("Error writing response: %v", err)
		}
	}

	log.Printf(
		"%s %s %s %d %s",
		r.Method,
		r.RequestURI,
		r.RemoteAddr,
		code,
		logErr.Error(),
	)
}
