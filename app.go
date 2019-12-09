package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func handler(response http.ResponseWriter, request *http.Request) {
	name, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(response, "This is v2 running in pod "+name, "\n")
}

func main() {
	fmt.Printf("Starting Web Server...")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
