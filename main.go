package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	// "os"
	// "strings"
)

type PayloadStruct struct {
	ProxyCondition string `json:"proxy_condition"`
}

var PORT = "80"
var URL_A = "http://example.com"
var URL_B = "https://rit.edu"
var URL_C = ""

// Because we love logggggggggg
func logSetup() {
	log.Printf("Server will run on: %s\n", PORT)
	log.Printf("URL A: %s\n", URL_A)
	log.Printf("URL B: %s\n", URL_B)
	log.Printf("URL C: %s\n", URL_C)

}

func requestBodyDecoder(request *http.Request) *json.Decoder {
	//Read body to buffer
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		panic(err)
	}
	
	// Because go lang is a pain in the ass if you read the body then any susequent calls
	// are unable to read the body again....
	request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer(body)))
}

func parseRequestBody(request *http.Request) PayloadStruct {
	decoder := requestBodyDecoder(request)

	var requestPayload PayloadStruct
	err := decoder.Decode(&requestPayload) 

	if err != nil {
		panic(err)
	}

	return requestPayload
}

func serveRequestProxy(target string, res http.ResponseWriter, req *http.Request) {
	// parse the url
	url, _ := url.Parse(target)

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Update the header
	req.URL.Host = url.Host
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host

	proxy.ServeHTTP(res, req)
}

func handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
	// We will get to this...
	requestPayload := parseRequestBody(req)
	log.Printf("Request payload: %s", requestPayload)

	var url = ""
	if requestPayload.ProxyCondition == "a" {
		// url a
		url = URL_A
	}else if requestPayload.ProxyCondition == "b" {
		// url b
		url = URL_B
	}else{
		// default
		url = URL_C
	}

	log.Printf("proxy_condition: %s, proxy_url: %s\n", 
				requestPayload.ProxyCondition, url)
	serveRequestProxy(url, res, req)

}

func main(){
	logSetup()

	// Start the server
	http.HandleFunc("/", handleRequestAndRedirect)
	log.Printf("Listening on port %s\n", PORT)
	if err := http.ListenAndServe(":" + PORT, nil); err != nil {
		panic(err)
	}

}