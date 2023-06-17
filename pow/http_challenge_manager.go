package pow

import (
	"net/http"
)

const (
	ChallengeHeader = "X-WOW-POW-CHALLENGE"
	SolutionHeader  = "X-WOW-POW-SOLUTION"
)

type HTTPChallengeManager struct {
}

func (c *HTTPChallengeManager) GetChallenge(response *http.Response) string {
	return response.Header.Get(ChallengeHeader)
}

func (c *HTTPChallengeManager) SetChallenge(response *http.Response) string {
	return response.Header.Get(ChallengeHeader)
}

func (c *HTTPChallengeManager) GetSolution(request *http.Request) string {
	return request.Header.Get(SolutionHeader)
}

func (c *HTTPChallengeManager) SetSolution(request *http.Request, solution string) {
	request.Header.Set(SolutionHeader, solution)
}
