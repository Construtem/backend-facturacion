package handlers

import (
	"backend-facturacion/models"
	"backend-facturacion/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// Simplificado para solo recibir el ID de cotización
type requestQuotepreview struct {
	CotizacionID int `json:"cotizacion_id"`
}

func CrearQuotePreview(c *gin.Context) {
	var request1 requestQuotepreview
	var preview1 models.QuotePreview

	c.ShouldBindJSON(&request1)

	// Buscar la cotización en la base de datos
	type Cotizacion struct {
		ID    int
		Total float64
	}
	var cotizacion Cotizacion
	db := utils.GetDB()
	if err := db.Table("cotizacions").Where("id = ?", request1.CotizacionID).First(&cotizacion).Error; err != nil {
		c.JSON(404, gin.H{"error": "Cotización no encontrada"})
		return
	}

	// Llenar los datos de QuotePreview
	preview1.CotizacionId = request1.CotizacionID
	preview1.Subtotal = cotizacion.Total
	preview1.Tax = preview1.Subtotal * 0.19
	preview1.Total = preview1.Subtotal + preview1.Tax
	preview1.PaymentStatus = models.Pending
	preview1.SuccessfulPaymentIntentID = "0"
	preview1.IssuedAt = time.Now()

	// Guardar en la base de datos
	if err := db.Create(&preview1).Error; err != nil {
		c.JSON(500, gin.H{"error": "No se pudo crear el preview"})
		return
	}

	// Usar la nueva función para crear la factura
	factura, err := utils.CrearFactura(request1.CotizacionID, int(preview1.ID))
	if err != nil {
		c.JSON(500, gin.H{"error": "No se pudo crear la factura: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"id de cotizacion": preview1.CotizacionId,
		"total":            preview1.Total,
		"id":               preview1.ID,
		"id_factura":       factura.ID,
	})
}
