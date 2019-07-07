package main

import (
	"fmt"
	"net/http"
	"os"
)

func handler(response http.ResponseWriter, request *http.Request) {
	name, _ := os.Hostname()
	fmt.Fprint(response, "This is v1 running in pod "+name, "\n")
}

func main() {
  fmt.Printf("Starting Web Server...")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

