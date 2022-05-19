package routes

import (
	"server/controllers"

	"github.com/gin-gonic/gin"
)

func OrderRoute(router *gin.Engine) {
	// All order routes come here.
	router.POST("/order/create", controllers.AddOrder())
	router.GET("/orders", controllers.GetAllOrders())
}
