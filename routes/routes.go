package routes

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter inicializa las rutas de la API
func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Aquí puedes agregar más rutas

	return r
}