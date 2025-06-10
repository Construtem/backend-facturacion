package routes

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter inicializa las rutas de la API
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Middleware CORS
	/*r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3050"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// Endpoint público
	r.GET("/api/cotizacion/:id", handlers.GetCotizacionByID)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})*/

	// Aquí se pueden agregar más rutas

	return r
}
