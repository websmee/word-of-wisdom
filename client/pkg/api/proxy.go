package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type Proxy struct {
	targetBaseURL string
}

func NewProxy(targetBaseURL string) *Proxy {
	return &Proxy{targetBaseURL}
}

func (s *Proxy) SendRequest(_ context.Context, request *http.Request) (*http.Response, error) {
	urlStr := s.targetBaseURL + request.RequestURI
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url '%s', err: %w", urlStr, err)
	}

	request.URL = u
	request.RequestURI = ""

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("http request failed, err: %w", err)
	}
	defer response.Body.Close()

	return response, nil
}
