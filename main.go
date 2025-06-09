package main

import (
	"backend-facturacion/routes"
	"backend-facturacion/utils"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Corriendo en localhost:8080 !!!")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error al cargar el archivo .env")
	}

	dsnFromEnv := os.Getenv("DATABASE_DSN")
	fmt.Println("DATABASE_DSN cargado por godotenv:", dsnFromEnv)

	utils.InitDB()
	router := routes.SetupRouter()
	router.Run() // listen and serve on 0.0.0.0:8080

}
