package http

import (
	"encoding/json"
	"net/http"
)

func WriteResponse(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if response != nil {
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			return
		}
	}
}

func WriteErrorResponse(w http.ResponseWriter, err interface{}) {
	if httpErr, ok := err.(HTTPError); ok {
		http.Error(w, httpErr.Message, httpErr.Code)
		return
	}
	if domainErr, ok := err.(error); ok {
		httpErr := MapDomainError(domainErr)
		http.Error(w, httpErr.Message, httpErr.Code)
		return
	}

	http.Error(w, ErrInternalServer.Message, ErrInternalServer.Code)
}
