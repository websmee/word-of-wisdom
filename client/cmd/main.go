package main

import (
	"os"

	"github.com/gin-gonic/gin"

	"websmee/word-of-wisdom/client/pkg/api"
	"websmee/word-of-wisdom/client/pkg/infrastructure"
	"websmee/word-of-wisdom/client/pkg/usecases"
)

var wowServiceBaseURL = os.Getenv("WOW_SERVICE_BASE_URL")

func main() {
	wowService := infrastructure.NewWOWService(wowServiceBaseURL)
	wowProxy := usecases.NewWOWProxy(wowService)
	wowProxyHandler := api.NewWOWRequestHandler(wowProxy)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/api/v1/wow/random", wowProxyHandler)

	_ = router.Run(":8081")
}
