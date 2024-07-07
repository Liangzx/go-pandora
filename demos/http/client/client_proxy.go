package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func main() {
	fmt.Println("hello")
	proxyURL, err := url.Parse("http://localhost:8888")
	if err != nil {
		fmt.Println("Error:", err)
	}
	client := &http.Client{
		CheckRedirect: nil,
		Transport: &http.Transport{
			// Proxy: http.ProxyFromEnvironment,
			Proxy: http.ProxyURL(proxyURL),
		},
	}
	// http_proxy="locahost:8080"
	response, err := client.Get("https://cloud.tencent.com/developer/article/1418457")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println(string(body))
}
