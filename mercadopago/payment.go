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
		fmt.Println(payment1.QuotePreviewID)

		// Actualizar el estado de pago de la quote preview
		db := utils.GetDB()
		err := db.Model(&models.QuotePreview{}).
			Where("id = ?", payment1.QuotePreviewID).
			Update("payment_status", payment1.Status).Error
		if err != nil {
			fmt.Println("Error actualizando estado de pago:", err)
		}

		// Obtener el ID de cotización desde el QuotePreview
		var quotePreview models.QuotePreview
		if err := db.Where("id = ?", payment1.QuotePreviewID).First(&quotePreview).Error; err != nil {
			fmt.Println("Error obteniendo quote preview:", err)
			c.JSON(500, gin.H{
				"message": "Error obteniendo información de cotización",
			})
			return
		}

		// Actualizar la columna pago_id en cotizaciones usando solo el ID de cotización
		var cotizacionID int
		errCot := db.Table("preview_cotizacion").
			Select("cotizacion_id").
			Where("id = ?", payment1.QuotePreviewID).
			Row().Scan(&cotizacionID)
		if errCot != nil {
			fmt.Println("Error obteniendo cotizacion_id desde quote_previews:", errCot)
		} else {
			if err := db.Table("cotizaciones").
				Where("id = ?", cotizacionID).
				Update("pago_id", payment1.PagoID).Error; err != nil {
				fmt.Println("Error actualizando pago_id en cotizaciones:", err)
			} else {
				fmt.Println("pago_id actualizado correctamente en cotizaciones")
			}
		}

		// Enviar el estado usando el ID de cotización real
		if err := utils.EstadoPago(payment1.Status, quotePreview.CotizacionId); err != nil {
			c.JSON(400, gin.H{
				"message": "No se pudo enviar el status de la cotizacion",
			})
		}

		if status == "approved" {
			// Actualizar el estado de la factura a "aprobado"
			if err := db.Model(&models.Factura{}).
				Where("quote_preview_id = ?", payment1.QuotePreviewID).
				Update("estado", "aprobado").Error; err != nil {
				fmt.Println("Error actualizando estado de factura:", err)
			} else {
				fmt.Println("Factura actualizada correctamente a estado 'aprobado'")
			}
		}

		if status == "rejected" {
			// Actualizar el estado de la factura a "rechazado"
			if err := db.Model(&models.Factura{}).
				Where("quote_preview_id = ?", payment1.QuotePreviewID).
				Update("estado", "rechazado").Error; err != nil {
				fmt.Println("Error actualizando estado de factura:", err)
			} else {
				fmt.Println("Factura actualizada correctamente a estado 'rechazado'")
			}
		}

		if status == "in_process" {
			// Actualizar el estado de la factura a "en proceso"
			if err := db.Model(&models.Factura{}).
				Where("quote_preview_id = ?", payment1.QuotePreviewID).
				Update("estado", "en proceso").Error; err != nil {
				fmt.Println("Error actualizando estado de factura:", err)
			} else {
				fmt.Println("Factura actualizada correctamente a estado 'en proceso'")
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
