package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	userInput := r.URL.Query().Get("userInput")
	w.Header().Set("X-Example", "static-value-"+userInput) // want "possible HTTP header injection found"

	w.Header().Set("X-Example", "static-value-"+userInput+"-suffix") // want "possible HTTP header injection found"
}
