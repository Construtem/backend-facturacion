package routes

import (
	"backend-facturacion/handlers"
	"backend-facturacion/mercadopago"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter inicializa las rutas de la API
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Middleware CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3050"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	//Endpoints
	r.POST("/API/v1/webhook", mercadopago.Webhook)
	r.POST("/API/v1/payment", mercadopago.Payment)
	r.GET("/api/cotizacion/:id", handlers.GetCotizacionByID)

	return r
}
