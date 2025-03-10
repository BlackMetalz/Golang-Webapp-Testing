package main

import (
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/stretchr/testify/assert"
)

func addSecurityHeaders(w http.ResponseWriter) {
    w.Header().Set("X-Content-Type-Options", "nosniff")
    w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
}

func logAndRespond(w http.ResponseWriter, statusCode int) {
    addSecurityHeaders(w)
    fmt.Printf("Returning status code: %d\n", statusCode)
    http.Error(w, http.StatusText(statusCode), statusCode)
}

func main() {
    // Handle root path specially
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // If it's not exactly "/", return 404
        if r.URL.Path != "/" {
            logAndRespond(w, http.StatusNotFound)
            return
        }
        addSecurityHeaders(w)
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "OK")
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

    http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Version: 4.0.0")
    })

    fmt.Println("Server starting on :8080")
    http.ListenAndServe(":8080", nil)
}

func TestEndpoints(t *testing.T) {
    // Test valid path
    t.Run("Valid path /", func(t *testing.T) {
        req := httptest.NewRequest("GET", "/", nil)
        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if r.URL.Path != "/" {
                logAndRespond(w, http.StatusNotFound)
                return
            }
            addSecurityHeaders(w)
            w.WriteHeader(http.StatusOK)
            fmt.Fprintf(w, "OK")
        })
        handler.ServeHTTP(rr, req)
        assert.Equal(t, http.StatusOK, rr.Code)
        assert.Equal(t, "nosniff", rr.Header().Get("X-Content-Type-Options"))
        assert.Equal(t, "same-origin", rr.Header().Get("Cross-Origin-Opener-Policy"))
    })

    // Test invalid path
    t.Run("Invalid path /abc", func(t *testing.T) {
        req := httptest.NewRequest("GET", "/abc", nil)
        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if r.URL.Path != "/" {
                logAndRespond(w, http.StatusNotFound)
                return
            }
            addSecurityHeaders(w)
            w.WriteHeader(http.StatusOK)
            fmt.Fprintf(w, "OK")
        })
        handler.ServeHTTP(rr, req)
        assert.Equal(t, http.StatusNotFound, rr.Code)
        assert.Equal(t, "nosniff", rr.Header().Get("X-Content-Type-Options"))
        assert.Equal(t, "same-origin", rr.Header().Get("Cross-Origin-Opener-Policy"))
    })
}
