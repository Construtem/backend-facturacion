package main

import (
	"backend-facturacion/mercadopago"

	"github.com/gin-gonic/gin"
)

func main() {

	mercadopago.Main()
	r := gin.Default()
	r.POST("/API/v1/webhook")
	//r.GET("/API/v1/hola", get)
	r.Run(":3050")

}
