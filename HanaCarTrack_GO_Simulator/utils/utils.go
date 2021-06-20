package utils

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Response struct {
	StatusCode int
	Headers    map[string][]string
	Response   []byte
}

func callService(urlString string, method string, headers map[string][]string, body string) Response {
	urlParse, _ := url.Parse(urlString)
	req := http.Request{
		Method: method,
		URL:    urlParse,
		Header: http.Header(headers),
		Body:   ioutil.NopCloser(strings.NewReader(body)),
	}
	res, err := http.DefaultClient.Do(&req)
	if err != nil {
		return Response{
			StatusCode: 500,
			Response:   []byte(err.Error()),
		}
	} else {
		resBytes, _ := ioutil.ReadAll(res.Body)
		return Response{
			StatusCode: res.StatusCode,
			Headers:    res.Header,
			Response:   resBytes,
		}
	}
}
