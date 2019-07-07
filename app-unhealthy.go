package main

import (
	"fmt"
	"net/http"
	"os"
)

var (
	counter = 0
)

func handler(response http.ResponseWriter, request *http.Request) {
	name, err := os.Hostname()
	checkError(err)
	if counter < 5 {
		fmt.Fprint(response, "You've hit "+name, "\n")
		counter++
	} else {
		response.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(response, "Some internal error has occurred! This is pod "+name, "\n")
	}

}

func main() {
	fmt.Printf("Starting Web Server...")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
