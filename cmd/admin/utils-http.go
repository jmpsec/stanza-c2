package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// ContentType for header key
const ContentType string = "Content-Type"

// JSONApplicationUTF8 for Content-Type headers
const JSONApplicationUTF8 string = "application/json; charset=UTF-8"

// Helper to send a serialized JSON response
func httpResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set(ContentType, JSONApplicationUTF8)
	content, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error serializing response: %v", err)
		content = []byte("error serializing response")
	}
	w.WriteHeader(code)
	_, _ = w.Write(content)
}
