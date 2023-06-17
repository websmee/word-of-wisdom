package infrastructure

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type WOWService struct {
	baseURL string
}

func NewWOWService(baseURL string) *WOWService {
	return &WOWService{baseURL}
}

func (s *WOWService) SendRequest(_ context.Context, request *http.Request) (*http.Response, error) {
	urlStr := s.baseURL + request.RequestURI
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

	return response, nil
}
