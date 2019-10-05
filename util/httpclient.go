package util

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/url"
)

var client = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
}

func DoHTTPRequest(method, endpoint, path string, headers, params map[string]string, body []byte) (*http.Response, error) {
	URL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	URL.Path += path
	parameters := url.Values{}
	for key, val := range params {
		parameters.Add(key, val)
	}

	URL.RawQuery = parameters.Encode()
	//fmt.Println("request URL is ", URL.String())
	req, err := http.NewRequest(method, URL.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	req.Close = true
	resp, err := client.Do(req)

	return resp, err
}

func ReleaseBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
}

func IsResponseStatusOk(resp *http.Response) bool {
	return http.StatusOK == resp.StatusCode/100*100
}

func CopyResponseBody(res *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(res.Body)
	//res.Body.Close()
	if err != nil {
		return nil, err
	}
	return body, nil
}
