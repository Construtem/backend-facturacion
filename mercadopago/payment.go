package mercadopago

import (
	"context"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/payment"
)

type Request struct {
	TransactionAmount float64 `json:"transaction_amount"`
	PaymentMethodID   string  `json:"payment_method_id"`
	Email             string  `json:"email"`
	CardToken         string  `json:"card_token"`
}

type PaymentStatus struct {
	Status             string `json:"status"`
	TransactionDetails struct {
		TotalPaidAmount float64 `json:"total_paid_amount"`
	} `json:"transaction_details"`
}

func Payment(c *gin.Context) {

	var request1 Request
	accessToken := accessToken()

	//se obtiene el card_token desde el front
	c.ShouldBind(&request1)

	//se crea un configuracion con el token de acceso
	cfg, err := config.New(accessToken)
	if err != nil {
		c.JSON(400, gin.H{

			"error": err,
		})
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
		return
	}

	//Analizar la respuesta de mercadopago, se guarda si fue exitoso o no y la cantidad total pagada
	var statusResp PaymentStatus
	bytes, _ := json.Marshal(resource)
	json.Unmarshal(bytes, &statusResp)

	//respuesta para el front del estado del pago
	c.JSON(200, gin.H{
		"status":            statusResp.Status,
		"total_paid_amount": statusResp.TransactionDetails.TotalPaidAmount,
	})

}
