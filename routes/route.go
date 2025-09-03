package routes

import (
	"blockchain/controllers"
	"os"

	"github.com/gin-gonic/gin"
)

func InitRoute() *gin.Engine {
	r := gin.Default()

	r.GET("/chain", controllers.GetChain)

	// transaction
	r.POST("/transaction", controllers.CreateTransaction)

	// block
	if os.Getenv("ENV") == "development" {
		r.POST("/block", controllers.ManuelMineBlock)
	}

	return r
}
