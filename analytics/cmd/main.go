package main

import (
    "log"
    "net/http"
)

func main() {
    // start a minimal HTTP server for search API (stub)
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("ok"))
    })
    log.Println("Analytics Service listening on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
