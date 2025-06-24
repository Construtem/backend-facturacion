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
		AllowOrigins:     []string{"https://facturacion.tssw.cl"},  // Para trabajar en local usar http://localhost:3000
		AllowMethods:     []string{"GET, POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	//Endpoints
	r.POST("/API/v1/webhook", mercadopago.Webhook)
	r.POST("/API/v1/payment", mercadopago.Payment)
	r.GET("/api/cotizacion/:id", handlers.GetCotizacionByID)
	r.POST("/api/preview_cotizacion", handlers.CreatePreviewCotizacion)

	return r
}
