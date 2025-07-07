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
	var statusResp models.PaymentResp
	bytes, _ := json.Marshal(resource)
	json.Unmarshal(bytes, &statusResp)

	var payment1 models.Payment_intent

	//Creacicon del payment_intent a travez de respuesta de mercadopago
	payment1.Status = statusResp.Status
	payment1.TransactionAmount = statusResp.TransactionDetails.TotalPaidAmount
	payment1.PagoID = statusResp.ID
	payment1.QuotePreviewID = request1.CotizacionID
	payment1.CreatedAt = statusResp.DateCreated
	payment1.MetodoPago = statusResp.PaymentMethodID

	// Guardar payment_intent en la base de datos
	db := utils.GetDB()
	if err := db.Create(&payment1).Error; err != nil {
		fmt.Println("Error guardando payment_intent:", err)
	}

	//verifiacion del status a la api de mercadopago
	status := utils.VerificarPago(payment1.PagoID)

	if status == payment1.Status {

		fmt.Println("Correcto")

		// Actualizar el estado de pago de la quote preview
		db := utils.GetDB()
		err := db.Model(&models.QuotePreview{}).
			Where("id = ?", payment1.QuotePreviewID).
			Update("payment_status", payment1.Status).Error
		if err != nil {
			fmt.Println("Error actualizando estado de pago:", err)
		}

		// Buscar la factura temporal
		var facturaTemp models.FacturaTemp
		if err := db.Where("quote_preview_id = ?", request1.CotizacionID).First(&facturaTemp).Error; err == nil {
			// Copiar los datos a una nueva factura real
			factura := models.Factura{
				CotizacionID:        facturaTemp.CotizacionID,
				QuotePreviewID:      facturaTemp.QuotePreviewID,
				RutCliente:          facturaTemp.RutCliente,
				TipoDocumento:       facturaTemp.TipoDocumento,
				Folio:               facturaTemp.Folio,
				FechaEmision:        facturaTemp.FechaEmision,
				FechaVencimiento:    facturaTemp.FechaVencimiento,
				TimbreElectronico:   facturaTemp.TimbreElectronico,
				SiiIndicacion:       facturaTemp.SiiIndicacion,
				FraseLegal:          facturaTemp.FraseLegal,
				RutEmisor:           facturaTemp.RutEmisor,
				RazonSocialEmisor:   facturaTemp.RazonSocialEmisor,
				GiroEmisor:          facturaTemp.GiroEmisor,
				DireccionEmisor:     facturaTemp.DireccionEmisor,
				ComunaEmisor:        facturaTemp.ComunaEmisor,
				CiudadEmisor:        facturaTemp.CiudadEmisor,
				TelefonoEmisor:      facturaTemp.TelefonoEmisor,
				EmailEmisor:         facturaTemp.EmailEmisor,
				RutReceptor:         facturaTemp.RutReceptor,
				RazonSocialReceptor: facturaTemp.RazonSocialReceptor,
				GiroReceptor:        facturaTemp.GiroReceptor,
				DireccionReceptor:   facturaTemp.DireccionReceptor,
				ComunaReceptor:      facturaTemp.ComunaReceptor,
				CiudadReceptor:      facturaTemp.CiudadReceptor,
				ContactoReceptor:    facturaTemp.ContactoReceptor,
				SubtotalNeto:        facturaTemp.SubtotalNeto,
				Iva19:               facturaTemp.Iva19,
				IvaRetenido:         facturaTemp.IvaRetenido,
				TotalFinal:          facturaTemp.TotalFinal,
				UrlPdf:              facturaTemp.UrlPdf,
				UrlVerificacion:     facturaTemp.UrlVerificacion,
				CreatedAt:           facturaTemp.CreatedAt,
				UpdatedAt:           facturaTemp.UpdatedAt,
			}
			// Guardar la factura real
			if err := db.Create(&factura).Error; err != nil {
				fmt.Println("Error guardando factura real:", err)
			} else {
				// Eliminar la factura temporal
				db.Delete(&facturaTemp)
			}
		}
	}

	//respuesta para el front del estado del pago
	c.JSON(200, gin.H{
		"id":                statusResp.ID,
		"status":            statusResp.Status,
		"total_paid_amount": statusResp.TransactionDetails.TotalPaidAmount,
		"cotizacion_id":     request1.CotizacionID,
		"detalle_status":    statusResp.StatusDetail,
	})

}
