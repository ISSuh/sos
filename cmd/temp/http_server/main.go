package httpserver

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	mux.
		log.Fatal(http.ListenAndServe(":9009", mux))
}
