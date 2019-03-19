package jike

import (
	"fmt"
	"io"
	"net/http"
)

var httpclient = &http.Client{}
var headers = map[string]string{
	"Origin":          "http://web.okjike.com",
	"Referer":         "http://web.okjike.com",
	"User-Agent":      "jikefm",
	"Accept":          "application/json",
	"Content-Type":    "application/json",
	"Accept-Language": "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7",
	"App-Version":     "5.3.0",
	"platform":        "web",
}

func client() *http.Client {
	return httpclient
}

func newRequest(method string, url string, body io.Reader) *http.Request {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		fmt.Println(err)
	}
	for key, value := range headers {
		request.Header.Set(key, value)
	}
	return request
}
