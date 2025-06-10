package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backend-facturacion/models"
	"backend-facturacion/utils" // para obtener la instancia DB

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
	var cotizacion models.Cotizacion
	result := DB.First(&cotizacion, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cotización no encontrada"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error en la base de datos"})
		}
		return
	}

	var items []map[string]interface{}
	err = json.Unmarshal([]byte(cotizacion.ItemsJSON), &items)
	if err != nil {
		items = []map[string]interface{}{}
	}

	c.JSON(http.StatusOK, gin.H{
		"id":            cotizacion.ID,
		"fecha_emision": cotizacion.FechaEmision,
		"subtotal":      cotizacion.Subtotal,
		"impuesto":      cotizacion.Impuesto,
		"total":         cotizacion.Total,
		"items":         items,
	})
}
