package musicapi

import (
	"fmt"
	"io"
	"net/http"
)

var httpclient = &http.Client{}
var headers = map[string]string{
	"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.119 Safari/537.36",
	"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
	"Accept-Language": "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7",
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
