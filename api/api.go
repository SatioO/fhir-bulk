package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/satioO/fhir/v2/domain"
)

func Error(w http.ResponseWriter, r *http.Request, err error, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err == nil {
		err = fmt.Errorf("nil err")
	}
	logErr := err

	errorMsgJSON, err := json.Marshal(domain.ErrorResponse{Message: err.Error()})

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

func SuccessJson(w http.ResponseWriter, r *http.Request, data interface{}) {
	jsonMsg, err := json.Marshal(data)
	if err != nil {
		Error(w, r, fmt.Errorf("serialising response failed: %w", err), 500)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		Success(w, r, jsonMsg)
	}
}

func Success(w http.ResponseWriter, r *http.Request, data []byte) {
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
