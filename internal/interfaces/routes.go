package interfaces

import "github.com/gin-gonic/gin"

func RegisterHTTPServer(chatRouter *ChatUseCase) *gin.Engine {
	r := gin.Default()
	r.Any("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/douyin/message/action", chatRouter.Send)

	return r
}
