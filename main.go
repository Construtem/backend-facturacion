package main

import (
	"backend-facturacion/mercadopago"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.POST("/API/v1/webhook")
	r.POST("/API/v1/payment", mercadopago.Payment)
	r.Run(":3050")

}
