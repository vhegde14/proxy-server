package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request: ", r)
	targetURL := r.URL
	// Attempting to construct proxy request from initial HTTP Request
	proxyRequest, err := http.NewRequest(r.Method, targetURL.String(), r.Body)
	if err != nil {
		http.Error(w, "Error creating Proxy Server Request: ", http.StatusInternalServerError)
		fmt.Println("Error creating Proxy Request")
		return
	}

	// Adding initial HTTP Request headers to Proxy Request headers
	for name, values := range r.Header {
		for _, value := range values {
			proxyRequest.Header.Add(name, value)
		}
	}

	// Sending Server Request
	response, err := http.DefaultTransport.RoundTrip(proxyRequest)
	if err != nil {
		http.Error(w, "Error sending Proxy Server Request: ", http.StatusInternalServerError)
		fmt.Println("Error sending Proxy Request")
		return
	}
	defer response.Body.Close()

	for name, values := range response.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	w.WriteHeader(response.StatusCode)
	io.Copy(w, response.Body)
}

func main() {
	// Creating Server on Port 8080
	server := &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(handleRequest),
	}

	// Attempting to start Proxy Server
	fmt.Println("Starting Proxy Server on Port 8080")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error attempting to start Proxy Server: ", err)
	}
}
