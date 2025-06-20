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

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error al cargar el archivo .env")
	}

	dsnFromEnv := os.Getenv("DATABASE_DSN")
	fmt.Println("DATABASE_DSN cargado por godotenv:", dsnFromEnv)

	utils.InitDB()

	r := routes.SetupRouter()
	r.Run(":3050")

}
