package main

import (
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/s1-service/api/v1/status", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK from s1-service"))
		return
	})

	http.HandleFunc("/s2-service/api/v1/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from s1-service"))
		return
	})
	log.Fatal(http.ListenAndServe(":8081", nil))
	log.Println("S1 HTTP server started on port 8081")
}
