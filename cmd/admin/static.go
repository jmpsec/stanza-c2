package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

)

// Handler to serve static content with the proper header
func staticHandler(w http.ResponseWriter, r *http.Request) {
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(requestDump))

	path := r.URL.Path
	if strings.HasSuffix(path, ".css") {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	} else if strings.HasSuffix(path, ".js") {
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	} else if strings.HasSuffix(path, ".eot") {
		w.Header().Set("Content-Type", "application/vnd.ms-fontobject")
	} else if strings.HasSuffix(path, ".svg") {
		w.Header().Set("Content-Type", "image/svg+xml")
	} else if strings.HasSuffix(path, ".ttf") {
		w.Header().Set("Content-Type", "application/x-font-ttf")
	} else if strings.HasSuffix(path, ".woff") {
		w.Header().Set("Content-Type", "application/font-woff")
	} else if strings.HasSuffix(path, ".woff2") {
		w.Header().Set("Content-Type", "application/font-woff2")
	} else if strings.HasSuffix(path, ".otf") {
		w.Header().Set("Content-Type", "application/x-font-otf")
	} else if strings.HasSuffix(path, ".ico") {
		w.Header().Set("Content-Type", "image/x-icon")
	} else if strings.HasSuffix(path, ".gif") {
		w.Header().Set("Content-Type", "image/gif")
	} else if strings.HasSuffix(path, ".png") {
		w.Header().Set("Content-Type", "image/png")
	}

	http.ServeFile(w, r, path)
}

// Handler for the favicon
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	http.ServeFile(w, r, "./static/img/favicon.png")
}
