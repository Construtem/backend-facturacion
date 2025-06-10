package main

import (
	"backend-facturacion/mercadopago"
	"time"

	"backend-facturacion/handlers"
	//"backend-facturacion/routes"
	"backend-facturacion/utils"
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	r := gin.Default()
	r.POST("/API/v1/webhook")
	r.POST("/API/v1/payment", mercadopago.Payment)
	r.GET("/api/cotizacion/:id", handlers.GetCotizacionByID)

	// Middleware CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// Endpoint público

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error al cargar el archivo .env")
	}

	dsnFromEnv := os.Getenv("DATABASE_DSN")
	fmt.Println("DATABASE_DSN cargado por godotenv:", dsnFromEnv)

	utils.InitDB()
	r.Run(":3050")

}
