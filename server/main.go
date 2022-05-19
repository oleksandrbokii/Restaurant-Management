package main

import (
	"os"
	"server/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := gin.Default()

	routes.OrderRoute(router)

	router.Run(":" + port)
}
