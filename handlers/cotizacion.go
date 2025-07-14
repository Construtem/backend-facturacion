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

func GetCotizacionByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	DB := utils.GetDB() // función para obtener la instancia DB
	if DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Conexión a la base de datos no inicializada correctamente."})
		return
	}

	var cotizacion models.QuotePreview
	result := DB.First(&cotizacion, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cotización no encontrada"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error en la base de datos"})
		}
		return
	}

	// Preparar la respuesta básica
	response := gin.H{
		"id":            cotizacion.ID,
		"fecha_emision": cotizacion.IssuedAt.Format(time.RFC3339),
		"subtotal":      cotizacion.Subtotal,
		"impuesto":      cotizacion.Tax,
		"total":         cotizacion.Total,
		"Cotizacion ID": cotizacion.CotizacionId,
	}

	// Obtener los datos del usuario usando el cotizacion_id del QuotePreview encontrado
	usuario, err := utils.ObtenerUsuario(cotizacion.CotizacionId)
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
