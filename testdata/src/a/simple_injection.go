package a

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	userSuppliedValue := r.URL.Query().Get("headerValue")

	// An attacker could supply a value containing newline (\n) or carriage return (\r) characters,
	// leading to an HTTP header injection vulnerability.
	w.Header().Set("X-Custom-Header", userSuppliedValue) // want "potential HTTP header injection"
	w.Write([]byte("Hello, world!"))
}
