package main

import (
	"fmt"
	"log"
	"net/http"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Got request. Method=%s. Path=%s\n", r.Method, r.URL.Path)
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	path := r.URL.Path
	fmt.Printf("%s\n", path)
}

func main() {
	http.HandleFunc("/update/", uploadHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
