package main

import (
	"backend-facturacion/routes"
	"fmt"
)

func main() {
	fmt.Println("Corriendo en localhost:8080 !!!")
	
	router := routes.SetupRouter()
  	router.Run() // listen and serve on 0.0.0.0:8080
	
}