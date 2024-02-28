package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

var inputFile = "./csv/input.csv"
var outputFile = "./csv/output.csv"

func main() {
	fmt.Println(" Starting csv check")

	// Read input file
	altData, _ := os.ReadFile(inputFile)
	
	// Split input file into array
	altDataArray := strings.Split(string(altData), "\r\n")

	// For Loop to check each entry
		// Check if old redirect is linking to new one -> If not set var to redirected URL, else set to "Correct Redirect"
		// Check response code of new redirect -> Save to Var

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	client2 := &http.Client{}
	

	for i := 1; i < len(altDataArray); i++ {
		fmt.Println("|-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-|")
		fmt.Println("Checking: " + altDataArray[i])

		var oldRedirect = strings.Split(altDataArray[i], ";")[0]
		var destination = strings.Split(altDataArray[i], ";")[1]
		
		var redirect = ""
		var statusCodeOrigin = 0
		var statusCodeDest = 0

		request, _ := http.NewRequest("GET", oldRedirect, nil)
		response, _ := client.Do(request)
		response2, _ := client2.Do(request)

		switch response.StatusCode {
		case 301:
			statusCodeOrigin = 301
			if response2.Request.URL.String() == destination {
			} else {
				redirect = response2.Request.URL.String()
			}
		case 404:
			statusCodeOrigin = 404
		case 200:
			statusCodeOrigin = 200
		default:
			statusCodeOrigin = response.StatusCode
		}

		response.Body.Close()

		request, _ = http.NewRequest("GET", destination, nil)
		response, _ = client.Do(request)

		statusCodeDest = response.StatusCode

		file, _ :=  os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		file.WriteString(oldRedirect + "; " + destination + "; " + fmt.Sprint(statusCodeOrigin) + "; " + fmt.Sprint(statusCodeDest) + "; " + redirect + "\r\n")
		
		fmt.Println("Checked Redirect")
		fmt.Println("|-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-/-|")
	}
}