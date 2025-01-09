package main

import (
    "fmt"
    "net/http"
)

func logAndRespond(w http.ResponseWriter, statusCode int) {
    fmt.Printf("Returning status code: %d\n", statusCode)
    http.Error(w, http.StatusText(statusCode), statusCode)
}

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        logAndRespond(w, http.StatusOK)
    })
	
    http.HandleFunc("/499", func(w http.ResponseWriter, r *http.Request) {
        logAndRespond(w, 499)
    })

    http.HandleFunc("/500", func(w http.ResponseWriter, r *http.Request) {
        logAndRespond(w, http.StatusInternalServerError)
    })

    http.HandleFunc("/502", func(w http.ResponseWriter, r *http.Request) {
        logAndRespond(w, http.StatusBadGateway)
    })

    http.HandleFunc("/503", func(w http.ResponseWriter, r *http.Request) {
        logAndRespond(w, http.StatusServiceUnavailable)
    })

    http.HandleFunc("/504", func(w http.ResponseWriter, r *http.Request) {
        logAndRespond(w, http.StatusGatewayTimeout)
    })

    // Start the server on port 8080
    http.ListenAndServe(":8080", nil)
}