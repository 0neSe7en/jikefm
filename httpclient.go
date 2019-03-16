package main

import (
	"fmt"
	"io"
	"net/http"
)

var client = &http.Client{}

func Client() *http.Client {
	return client
}

func NewRequest(method string, url string, body io.Reader) *http.Request {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		fmt.Println(err)
	}
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	request.Header.Set("App-Version", "5.3.0")
	request.Header.Set("platform", "web")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.119 Safari/537.36")
	request.Header.Set("Accept-Language", "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7")
	//request.Header.Set("x-jike-access-token", "")
	return request
}
