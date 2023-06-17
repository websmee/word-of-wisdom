package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"websmee/word-of-wisdom/server/pkg/usecases"
)

func NewWOWReadHandler(wowReader *usecases.WOWReader) func(c *gin.Context) {
	return func(c *gin.Context) {
		wow, err := wowReader.ReadRandomWOW(c)
		if err != nil {
			c.Error(fmt.Errorf("read word of wisdom handler failed, err: %w", err))
			c.Status(http.StatusInternalServerError)
			return
		}

		c.IndentedJSON(http.StatusOK, wow)
	}
}
