package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// KeepAlive for Connection headers
const KeepAlive string = "Keep-Alive"

// ContentType for header key
const ContentType string = "Content-Type"

// ContentDescription for header key
const ContentDescription string = "Content-Description"

// ContentDisposition for header key
const ContentDisposition string = "Content-Disposition"

// ContentLength for header key
const ContentLength string = "Content-Length"

// Connection for header key
const Connection string = "Connection"

// ContentTransferEncoding for header key
const ContentTransferEncoding string = "Content-Transfer-Encoding"

// Expires for header key
const Expires string = "Expires"

// CacheControl for header key
const CacheControl string = "Cache-Control"

// Pragma for header key
const Pragma string = "Pragma"

// PragmaPublic for header key
const PragmaPublic string = "public"

// TransferEncodingBinary for header key
const TransferEncodingBinary string = "binary"

// CacheControlMustRevalidate for header key
const CacheControlMustRevalidate string = "must-revalidate, post-check=0, pre-check=0"

// OctetStream for Content-Type headers
const OctetStream string = "application/octet-stream"

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
