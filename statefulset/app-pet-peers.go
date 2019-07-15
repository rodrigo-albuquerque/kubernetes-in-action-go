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

func checkError(err error, message string) {
	if err != nil {
		log.Println(message)
		log.Fatal(err)
	}
}

func sliceContains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func getHostnames(serviceName string) []string {
	// Receives an SRV record and Returns list of hostnames
	_, hostNames, err := net.LookupSRV("", "", serviceName)
	if err != nil {
		log.Fatal(err)
	}
	var hosts []string
	for _, host := range hostNames {
		hosts = append(hosts, host.Target)
	}
	return hosts
}

func getIPs(srvs []string, excludeLocalIP bool) []string {
	// excludeLocalIP means local POD's IP will be excluded when true
	var ips []string
	for _, srv := range srvs {
		if excludeLocalIP == false {
			listOfIps, err := net.LookupIP(srv)
			if err != nil {
				log.Println("Failed lookup")
				log.Fatal(err)
			}
			for _, ip := range listOfIps {
				ips = append(ips, ip.String())
			}
		} else {
			name, err := os.Hostname()
			checkError(err, "Cannot retrieve hostname from OS")
			if name != srv[:len(name)] {
				listOfIps, err := net.LookupIP(srv)
				if err != nil {
					log.Println("Failed lookup")
					log.Fatal(err)
				}
				for _, ip := range listOfIps {
					ips = append(ips, ip.String())
				}
			}
		}
	}
	return ips
}

func getRequestPod(url string) string {
	// Makes GET request to URL and returns only reply's body
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	requestBody, err := ioutil.ReadAll(resp.Body)
	checkError(err, "Cannot not retrieve HTTP body message")
	return string(requestBody)
}

func myhandler(resp http.ResponseWriter, req *http.Request) {
	name, err := os.Hostname()
	checkError(err, "Cannot retrieve hostname from OS")
	// POST requests are supposed to store data locally in /var/data/kubia.txt
	if req.Method == "POST" {
		postBody, err := ioutil.ReadAll(req.Body)
		checkError(err, "Cannot not retrieve HTTP body message")
		// If file doesn't exist, create it, or append to file
		f, err := os.OpenFile("/var/data/kubia.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		checkError(err, "Cannot not open file")
		defer f.Close()
		n, err := f.WriteString(string(postBody) + "\n")
		fmt.Printf("Wrote %d bytes\n", n)
		fmt.Fprint(resp, "Data stored on pod "+name, "\n")
	} else {
		// If GET request, we return our data
		fmt.Fprint(resp, "You've hit "+name, "\n")
		d, err := ioutil.ReadFile("/var/data/kubia.txt")
		if err != nil {
			fmt.Fprint(resp, "Data stored on this pod: No data posted yet\n")
		} else {
			fmt.Fprint(resp, "Data stored on this pod: \n")
			fmt.Fprint(resp, string(d))
		}
		hosts := getHostnames("kubia")
		ips := getIPs(hosts, false)
		s := strings.Split(req.RemoteAddr, ":")
		remoteIP, _ := s[0], s[1]
		// This part is only executed if source address is NOT from one of the other pods
		// i.e. don't request data from other pods if request comes from one of our pods
		if sliceContains(ips, remoteIP) == false {
			remoteIPs := getIPs(hosts, true)
			for _, ip := range remoteIPs {
				url := "http://" + ip + ":8080"
				reqBody := getRequestPod(url)
				fmt.Fprint(resp, reqBody)
			}
		}
	}
}

func main() {
	fmt.Println("Web server started on port 8080...")
	http.HandleFunc("/", myhandler)
	http.ListenAndServe(":8080", nil)
}

