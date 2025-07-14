package handlers

import (
	"backend-facturacion/models"
	"backend-facturacion/utils" // para obtener la instancia DB
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func PostPago(c *gin.Context) {
	// Obtener el ID del parámetro de la URL
	idStr := c.Param("id")
	quotePreviewID, err := strconv.Atoi(idStr)
	if err != nil || quotePreviewID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Obtener la instancia de la base de datos
	DB := utils.GetDB()
	if DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Conexión a la base de datos no inicializada correctamente."})
		return
	}

	// Buscar la factura por quote_preview_id
	var factura models.Factura
	result := DB.Where("quote_preview_id = ?", quotePreviewID).First(&factura)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Factura no encontrada para este preview de cotización"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error en la base de datos"})
		}
		return
	}

	// Preparar la respuesta con los datos solicitados
	response := gin.H{
		"numero_factura": factura.ID,                  // ID de la factura
		"nombre_cliente": factura.RazonSocialReceptor, // Nombre del cliente
		"empresa":        factura.RazonSocialEmisor,   // Empresa emisora
		"rut_cliente":    factura.RutCliente,          // RUT del cliente
	}

	c.JSON(200, response)
}
