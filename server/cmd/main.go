package main

import (
	"github.com/gin-gonic/gin"

	"websmee/word-of-wisdom/server/pkg/api"
	"websmee/word-of-wisdom/server/pkg/infrastructure"
	"websmee/word-of-wisdom/server/pkg/usecases"
)

func main() {
	wowRepository := infrastructure.NewWOWRepository([]string{
		"If you want to achieve greatness stop asking for permission. ~Anonymous",
		"Things work out best for those who make the best of how things work out. ~John Wooden",
		"To live a creative life, we must lose our fear of being wrong. ~Anonymous",
		"If you are not willing to risk the usual you will have to settle for the ordinary. ~Jim Rohn",
		"Trust because you are willing to accept the risk, not because it's safe or certain. ~Anonymous",
		"Take up one idea. Make that one idea your life - think of it, dream of it, live on that idea. Let the brain, muscles, nerves, every part of your body, be full of that idea, and just leave every other idea alone. This is the way to success. ~Swami Vivekananda",
		"All our dreams can come true if we have the courage to pursue them. ~Walt Disney",
		"Good things come to people who wait, but better things come to those who go out and get them. ~Anonymous",
		"If you do what you always did, you will get what you always got. ~Anonymous",
		"Success is walking from failure to failure with no loss of enthusiasm. ~Winston Churchill",
	})
	wowReader := usecases.NewWOWReader(wowRepository)
	wowReadHandler := api.NewWOWReadHandler(wowReader)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/api/v1/wow/random", wowReadHandler)

	_ = router.Run(":8080")
}
