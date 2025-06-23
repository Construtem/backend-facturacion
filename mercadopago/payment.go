package mercadopago

import (
	"context"
	"encoding/json"

	"fmt"

	"backend-facturacion/models"
	"backend-facturacion/utils"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/payment"
)

func Payment(c *gin.Context) {

	var request1 models.Request
	accessToken := os.Getenv("ACCESS_TOKEN")

	//se obtiene el card_token desde el front
	c.ShouldBind(&request1)

	//se crea un configuracion con el token de acceso
	cfg, err := config.New(accessToken)
	if err != nil {
		c.JSON(400, gin.H{

			"error": err,
		})

		fmt.Println("Error de token", err)
		return

	}

	client := payment.NewClient(cfg)

	//crear el pago
	request := payment.Request{
		TransactionAmount: request1.TransactionAmount,
		PaymentMethodID:   request1.PaymentMethodID,
		Payer: &payment.PayerRequest{
			Email: request1.Email,
		},
		Token:        request1.CardToken,
		Installments: 1,
	}

	//respuesta de mercadopago
	resource, err := client.Create(context.Background(), request)
	if err != nil {

		c.JSON(400, gin.H{

			"error": err,
		})
		fmt.Println("Error de respuesta mercadopago", err)
		return
	}

	//Analizar la respuesta de mercadopago, se guarda el "id", el "status" y la cantidad total pagada
	var statusResp models.PaymentStatus
	bytes, _ := json.Marshal(resource)
	json.Unmarshal(bytes, &statusResp)

	var payment1 models.Payment_intent

	//Creacicon del payment_intent a travez de respuesta de mercadopago
	payment1.Status = statusResp.Status
	payment1.TransactionAmount = statusResp.TransactionDetails.TotalPaidAmount
	payment1.PagoID = statusResp.ID
	payment1.QuotePreviewID = request1.CotizacionID
	fmt.Println(payment1)

	//verifiacion del status a la api de mercadopago
	status := utils.VerificarPago(payment1.PagoID)
	fmt.Println("este es el status de la api xd", status)

	if status == payment1.Status {

		fmt.Println("Correcto")
	}

	//respuesta para el front del estado del pago
	c.JSON(200, gin.H{
		"id":                statusResp.ID,
		"status":            statusResp.Status,
		"total_paid_amount": statusResp.TransactionDetails.TotalPaidAmount,
		"cotizacion_id":     request1.CotizacionID,
	})

}
