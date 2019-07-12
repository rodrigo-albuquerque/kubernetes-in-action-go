package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func myhandler(resp http.ResponseWriter, req *http.Request) {
	name, _ := os.Hostname()
	if req.Method == "POST" {
		postBody, err := ioutil.ReadAll(req.Body)
		checkError(err)
		err = ioutil.WriteFile("kubia.txt", []byte(postBody), 0644)
		checkError(err)
		fmt.Fprint(resp, "Data stored on pod "+name, "\n")
	} else {
		fmt.Fprint(resp, "You've hit "+name, "\n")
	}
}

func main() {
	http.HandleFunc("/", myhandler)
	http.ListenAndServe(":8080", nil)
}
