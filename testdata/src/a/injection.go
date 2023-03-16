package main

import (
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	userInput := r.URL.Query().Get("userInput")
	
	// detect patterns
	w.Header().Set("X-Example", "static-value-"+userInput) // want "possible HTTP header injection found"
	w.Header().Set("X-Example", "static-value-"+userInput+"-suffix") // want "possible HTTP header injection found"
	w.Header().Set("X-Example-2", strings.Replace("static-value", "-", userInput, -1)) // want "possible HTTP header injection found"
	
	// ok patterns
	w.Header().Set("X-Example", "static-value")
}
