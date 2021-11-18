package main

import (
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/s2-service/api/v1/status", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK from s2-service"))
		return
	})

	http.HandleFunc("/s2-service/api/v1/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from s2-service"))
		return
	})

	log.Fatal(http.ListenAndServe(":8082", nil))
	log.Println("S2 HTTP server started on port 8082")
}
