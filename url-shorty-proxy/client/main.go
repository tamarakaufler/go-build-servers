package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func main() {

	proxyURL := "http://localhost:7000"
	proxy, err := url.Parse(proxyURL)
	if err != nil {
		log.Println(err)
	}

	//creating the URL to be loaded through the proxy
	urlStr := "http://httpbin.org/get"
	url, err := url.Parse(urlStr)
	if err != nil {
		log.Println(err)
	}

	//adding the proxy settings to the Transport object
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxy),
	}

	//adding the Transport object to the http Client
	client := &http.Client{
		Transport: transport,
	}

	//generating the HTTP GET request
	request, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		log.Println(err)
	}

	//calling the URL
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
	}

	//getting the response
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}
	//printing the response
	log.Println(string(data))
}
