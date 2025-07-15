package routes

import (
	"backend-facturacion/handlers"
	"backend-facturacion/mercadopago"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter inicializa las rutas de la API
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Middleware de CORS
	// Configuración de CORS para permitir solicitudes desde los frontends especificados
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{ // Lista de URLs permitidas para CORS
			os.Getenv("FRONT_VENTAS_URL"),      // URL del frontend de ventas
			os.Getenv("FRONT_INVENTARIO_URL"),  // URL del frontend de inventario
			os.Getenv("FRONT_FACTURACION_URL"), // URL del frontend de facturación
		},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour, // Tiempo máximo de caché para las solicitudes CORS
	}))

	//Endpoints
	r.POST("/API/v1/webhook", mercadopago.Webhook)
	r.POST("/API/v1/payment", mercadopago.Payment)
	r.GET("/api/cotizacion/:id", handlers.GetCotizacionByID)
	r.POST("API/crear_preview", handlers.CrearQuotePreview)
	r.GET("/api/pdf/factura/:id", handlers.GenerateInvoicePDFHandler)
	r.GET("/API/v1/post-pago/:id", handlers.PostPago)

	//endpoint de prueba
	//r.GET("/api/cotizaciones/checkout/:id", handlers.PruebaVentas)

	return r
}
