package handlers

import (
	"backend-facturacion/models"
	"backend-facturacion/utils"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
)

type requestQuotepreview struct {
	CotizacionID      int             `json:"cotizacion_id"`
	RutCliente        string          `json:"rut_cliente"`
	TipoDocumento     string          `json:"tipo_documento"`
	RutEmisor         string          `json:"rut_emisor"`
	EmailEmisor       string          `json:"email_emisor"`
	RutReceptor       string          `json:"rut_receptor"`
	DireccionReceptor string          `json:"direccion_receptor"`
	ComunaReceptor    string          `json:"comuna_receptor"`
	CiudadReceptor    string          `json:"ciudad_receptor"`
	ContactoReceptor  string          `json:"contacto_receptor"`
	Items             json.RawMessage `json:"items"`
}

func CrearQuotePreview(c *gin.Context) {
	var request1 requestQuotepreview
	var preview1 models.QuotePreview
	var factura1 models.FacturaTemp
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

	// Llenar los datos de Factura
	factura1.CotizacionID = request1.CotizacionID
	factura1.QuotePreviewID = int(preview1.ID)
	factura1.RutCliente = request1.RutCliente
	factura1.TipoDocumento = request1.TipoDocumento
	factura1.FechaEmision = time.Now()
	factura1.FechaVencimiento = time.Now()
	factura1.RutEmisor = request1.RutEmisor
	factura1.EmailEmisor = request1.EmailEmisor
	factura1.RutReceptor = request1.RutReceptor
	factura1.DireccionReceptor = request1.DireccionReceptor
	factura1.ComunaReceptor = request1.ComunaReceptor
	factura1.CiudadReceptor = request1.CiudadReceptor
	factura1.ContactoReceptor = request1.ContactoReceptor
	factura1.SubtotalNeto = preview1.Subtotal
	factura1.TotalFinal = preview1.Total
	factura1.Items = request1.Items

	// Campos constantes o inventados
	factura1.Folio = "FOLIO123"
	factura1.TimbreElectronico = "TIMBRE"
	factura1.SiiIndicacion = "SII"
	factura1.FraseLegal = "Frase legal"
	factura1.RazonSocialEmisor = "Empresa S.A."
	factura1.GiroEmisor = "Servicios"
	factura1.DireccionEmisor = "Calle Falsa 123"
	factura1.ComunaEmisor = "Comuna"
	factura1.CiudadEmisor = "Ciudad"
	factura1.TelefonoEmisor = "123456789"
	factura1.RazonSocialReceptor = "Cliente S.A."
	factura1.GiroReceptor = "Comercio"
	factura1.Iva19 = factura1.SubtotalNeto * 0.19
	factura1.IvaRetenido = 0
	factura1.UrlPdf = ""
	factura1.UrlVerificacion = ""
	factura1.CreatedAt = time.Now()
	factura1.UpdatedAt = time.Now()

	// Guardar factura en la base de datos
	if err := db.Create(&factura1).Error; err != nil {
		c.JSON(500, gin.H{"error": "No se pudo crear la factura"})
		return
	}

	c.JSON(200, gin.H{
		"id de cotizacion": preview1.CotizacionId,
		"total":            preview1.Total,
		"id":               preview1.ID,
		"id_factura":       factura1.ID,
	})
}
