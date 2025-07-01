package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.String())
		w.Write([]byte("Hello, World!"))
	})
	http.ListenAndServe(":8080", nil)
}
