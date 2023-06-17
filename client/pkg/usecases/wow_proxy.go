package usecases

import (
	"context"
	"fmt"
	"net/http"
)

type WOWService interface {
	SendRequest(ctx context.Context, request *http.Request) (*http.Response, error)
}

type POWChallengeSolver interface {
	IsChallenge(ctx context.Context, response *http.Response) bool
	SolveChallenge(ctx context.Context, response *http.Response) (*http.Request, error)
}

type WOWProxy struct {
	wowService         WOWService
	powChallengeSolver POWChallengeSolver
}

func NewWOWProxy(wowService WOWService, powChallengeSolver POWChallengeSolver) *WOWProxy {
	return &WOWProxy{wowService, powChallengeSolver}
}

func (p *WOWProxy) ProxyRequest(ctx context.Context, request *http.Request) (*http.Response, error) {
	response, err := p.wowService.SendRequest(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("word of wisdom request failed, err: %w", err)
	}

	if p.powChallengeSolver.IsChallenge(ctx, response) {
		request, err = p.powChallengeSolver.SolveChallenge(ctx, response)
		if err != nil {
			return nil, fmt.Errorf("failed to solve proof of work challenge, err: %w", err)
		}

		response, err = p.wowService.SendRequest(ctx, request)
		if err != nil {
			return nil, fmt.Errorf("word of wisdom request failed, err: %w", err)
		}
	}

	return response, nil
}
