package mercadopago

import (
	"context"
	"fmt"

	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/payment"
)

func Main() {
	accessToken := accessToken()

	cfg, err := config.New(accessToken)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := payment.NewClient(cfg)

	request := payment.Request{
		TransactionAmount: 1000,
		PaymentMethodID:   "visa",
		Payer: &payment.PayerRequest{
			Email: "prueba@prueba.com",
		},
		Token:        "d49464ff30484565d36230691a102bb0",
		Installments: 1,
	}

	resource, err := client.Create(context.Background(), request)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resource)
}
