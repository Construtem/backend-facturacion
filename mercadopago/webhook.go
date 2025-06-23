package mercadopago

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

var request map[string]interface{}

func Webhook(c *gin.Context) {

	c.BindJSON(&request)

	fmt.Println(request)

	c.JSON(200, gin.H{

		"message":   "Webhook Recibido",
		"respuesta": request,
	})
}
