package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func retrievePodSRVs(serviceName string) []*net.SRV {
	_, srvs, err := net.LookupSRV("", "", serviceName)
	if err != nil {
		log.Fatal(err)
	}
	return srvs
}

func getPodData(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return string(body)
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
		srvs := retrievePodSRVs("kubia")
		name, err := os.Hostname()
		checkError(err)
		for _, srv := range srvs {
			// If service is local, just respond directly
			if strings.Contains(srv.Target, string(name)) {
				// reply with pod's hostname
				fmt.Fprint(resp, "You've hit "+name, "\n")
				fmt.Fprint(resp, "Data stored on this pod: ")
				// If kubia.txt exists, fwd its data, otherwise just say there is nothing in it
				d, err := ioutil.ReadFile("/var/data/kubia.txt")
				if err != nil {
					fmt.Fprint(resp, "Data stored on this pod: No data posted yet\n")
				} else {
					fmt.Fprint(resp, string(d))
				}
				// Otherwise, if pod is remote, issue a GET request and print response
			} else {
				respBody := getPodData(srv.Target)
				fmt.Fprint(resp, respBody)
			}
		}
	}
}

func main() {
	http.HandleFunc("/", myhandler)
	http.ListenAndServe(":8080", nil)
}

