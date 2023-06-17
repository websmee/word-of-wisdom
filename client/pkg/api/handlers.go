package api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/websmee/word-of-wisdom/pow"
)

func NewPOWProxyHandler(proxy *Proxy, powManager *pow.HTTPChallengeManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		response, err := proxy.SendRequest(c, c.Request)
		if err != nil {
			c.Error(fmt.Errorf("initial request failed, err: %w", err))
			c.Status(http.StatusInternalServerError)
			return
		}

		data, _ := io.ReadAll(response.Body)
		c.Writer.Write(data)
	}
}
