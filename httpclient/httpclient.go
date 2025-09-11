package httpclient

import (
	"io"
	"net/http"
	"time"
)

// Status codes que merecem retry
var retryableStatus = map[int]bool{
	http.StatusRequestTimeout:      true, // 408
	http.StatusTooManyRequests:     true, // 429
	http.StatusInternalServerError: true, // 500
	http.StatusBadGateway:          true, // 502
	http.StatusServiceUnavailable:  true, // 503
	http.StatusGatewayTimeout:      true, // 504
}

type OptionsHttpclient struct {
	RetryCount int
	Timeout    int
}

type HttpClient struct {
	url        string
	headers    map[string]string
	retryCount int
	timeout    int
}

func NewHttpClient(ops OptionsHttpclient) *HttpClient {
	header := make(map[string]string)
	header["Content-Type"] = "application/json"

	return &HttpClient{
		headers:    header,
		retryCount: ops.RetryCount,
		timeout:    ops.Timeout,
	}
}

func (h *HttpClient) SetUrl(url string) *HttpClient {
	h.url = url
	return h
}

func (h *HttpClient) SetHeader(key, value string) *HttpClient {
	h.headers[key] = value
	return h
}

func (h *HttpClient) SetBearerToken(token string) *HttpClient {
	h.headers["Autorization"] = "Bearer " + token
	return h
}

func (h *HttpClient) SendGet() ([]byte, int, error) {
	req, err := http.NewRequest("GET", h.url, nil)

	if err != nil {
		return nil, 500, err
	}

	setHeaderInNewRequest(h.headers, req)

	response, statusCode, err := sendClient(h, req)

	if err != nil {
		return nil, statusCode, err
	}

	return response, statusCode, err
}

func (h *HttpClient) SendPost() {

}

func setHeaderInNewRequest(headers map[string]string, h *http.Request) {
	for key, value := range headers {
		h.Header.Set(key, value)
	}
}

func sendClient(h *HttpClient, request *http.Request) ([]byte, int, error) {

	client := &http.Client{Timeout: time.Duration(h.timeout) * time.Second}

	resp, err := client.Do(request)
	if err != nil {
		return nil, 500, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return bodyBytes, resp.StatusCode, nil
}
