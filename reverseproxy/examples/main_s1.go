/*
	service S1 must be started before main.go: the reverse proxy is going to proxy the requests prefixed with "s1" to this service 
*/
package main

import (
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/api/v1/status", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK from s1-service"))
		return
	})

	http.HandleFunc("/api/v1/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from s1-service"))
		return
	})
	
	log.Println("S1 HTTP server started on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
	
}
