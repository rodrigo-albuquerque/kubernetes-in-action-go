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
	name, err := os.Hostname()
	checkError(err)
	if req.Method == "POST" {
		postBody, err := ioutil.ReadAll(req.Body)
		checkError(err)
		// If file doesn't exist, create it, or append to file
		f, err := os.OpenFile("/var/data/kubia.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		checkError(err)
		defer f.Close()
		n, err := f.WriteString(string(postBody) + "\n")
		fmt.Printf("Wrote %d bytes\n", n)
		fmt.Fprint(resp, "Data stored on pod "+name, "\n")
	} else {
		fmt.Fprint(resp, "You've hit "+name, "\n")
		fmt.Fprint(resp, "Data stored on this pod: ")
		d, err := ioutil.ReadFile("/var/data/kubia.txt")
		if err != nil {
			fmt.Fprint(resp, "Data stored on this pod: No data posted yet\n")
		} else {
			fmt.Fprint(resp, string(d))
		}
	}
}

func main() {
	http.HandleFunc("/", myhandler)
	http.ListenAndServe(":8080", nil)
}
