package handlers

import (
	"backend-facturacion/models"
	"backend-facturacion/utils" // para obtener la instancia DB
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type requestQuotepreview struct {
	CotizacionID int `json:"cotizacion_id"`
}
type Cotizacion struct {
	ID    int
	Total float64
}

func GetCotizacionByID(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(400, gin.H{"error": "ID inválido"})
		return
	}

	DB := utils.GetDB() // función para obtener la instancia DB
	if DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Conexión a la base de datos no inicializada correctamente."})
		return
	}

	var request1 requestQuotepreview
	var preview1 models.QuotePreview

	// Asignar el ID de llegada a request1
	request1.CotizacionID = id

	var cotizacion Cotizacion
	db := utils.GetDB()
	if err := db.Table("cotizaciones").Where("id = ?", request1.CotizacionID).First(&cotizacion).Error; err != nil {
		c.JSON(404, gin.H{"error": "Cotización no encontrada"})
		return
	}

	// Verificar si ya existe un QuotePreview para esta cotización
	var statusPagado models.PaymentStatus
	err = db.Where("cotizacion_id = ?", request1.CotizacionID).First(&preview1).Error
	if err != nil {
		// Si el error es que no se encontró el registro, crear un nuevo QuotePreview
		if err == gorm.ErrRecordNotFound {
			// Llenar los datos de QuotePreview
			preview1.CotizacionId = request1.CotizacionID
			preview1.Subtotal = cotizacion.Total / 1.19
			preview1.Tax = preview1.Subtotal * 0.19
			preview1.Total = preview1.Subtotal + preview1.Tax
			preview1.PaymentStatus = models.Pending
			preview1.IssuedAt = time.Now()

			// Guardar en la base de datos
			if err := db.Create(&preview1).Error; err != nil {
				c.JSON(500, gin.H{"error": "No se pudo crear el preview"})
				return
			}

			// Como es nuevo, no está pagado
			statusPagado = "nuevo"
		} else {
			// Si es otro tipo de error de base de datos
			c.JSON(500, gin.H{"error": "Error al buscar el quote preview: " + err.Error()})
			return
		}
	} else {
		// Si encontramos el QuotePreview, verificar si está pagado
		// Verificar el estado de pago basándose en PaymentStatus
		statusPagado = preview1.PaymentStatus
	}

	// Verificar si ya existe una factura para esta cotización
	var factura models.Factura
	err = db.Where("cotizacion_id = ?", request1.CotizacionID).First(&factura).Error
	if err != nil {
		// Si el error es que no se encontró el registro, crear una nueva factura
		if err == gorm.ErrRecordNotFound {
			facturaPtr, errCrear := utils.CrearFactura(request1.CotizacionID, int(preview1.ID))
			if errCrear != nil {
				c.JSON(500, gin.H{"error": "No se pudo crear la factura: " + errCrear.Error()})
				return
			}
			factura = *facturaPtr
		} else {
			// Si es otro tipo de error de base de datos
			c.JSON(500, gin.H{"error": "Error al buscar la factura: " + err.Error()})
			return
		}
	}

	// Preparar la respuesta básica
	response := gin.H{
		"id":            preview1.ID,
		"fecha_emision": preview1.IssuedAt.Format(time.RFC3339),
		"subtotal":      preview1.Subtotal,
		"impuesto":      preview1.Tax,
		"total":         preview1.Total,
		"cotizacion_id": preview1.CotizacionId,
		"factura_id":    factura.ID,
		"pagado":        statusPagado,
	}

	// Obtener los datos del usuario usando el cotizacion_id del QuotePreview encontrado
	usuario, err := utils.ObtenerUsuario(preview1.CotizacionId)
	if err != nil {
		// Si hay error obteniendo el usuario, log el error pero continúa con la respuesta básica
		// No falles toda la petición por esto
		fmt.Printf("Error obteniendo datos del usuario: %v\n", err)
	} else {
		// Agregar los datos del usuario a la respuesta
		response["usuario"] = gin.H{
			"nombre": usuario.Nombre,
			"email":  usuario.Email,
		}
	}

	c.JSON(http.StatusOK, response)
}
