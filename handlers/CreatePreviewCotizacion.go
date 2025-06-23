package handlers

import (
	"log"
	"net/http"
	"time"

	"backend-facturacion/models"
	"backend-facturacion/utils"

	"github.com/gin-gonic/gin"
)

type CreateQuotePreviewRequest struct {
	FechaEmision string  `json:"fecha_emision"`
	Subtotal     float64 `json:"subtotal"`
	Impuesto     float64 `json:"impuesto"`
	Total        float64 `json:"total"`
}

func CreatePreviewCotizacion(c *gin.Context) {
	var req CreateQuotePreviewRequest

	// Validación del JSON de entrada
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON de entrada inválido o malformado. Verifique el formato."})
		return
	}

	// Validación de los valores numéricos
	if req.Subtotal < 0 || req.Impuesto < 0 || req.Total < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Subtotal, impuesto y total no pueden ser valores negativos."})
		return
	}

	if req.Subtotal+req.Impuesto != req.Total {
		c.JSON(http.StatusBadRequest, gin.H{"error": "La suma de subtotal e impuesto no concuerda con el total. Verifique los montos."})
		return
	}

	// Validación de la fecha
	fecha, err := time.Parse(time.RFC3339, req.FechaEmision)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido (debe ser RFC3339)"})
		return
	}
	if fecha.After(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "La fecha de emisión no puede ser en el futuro"})
		return
	}

	// Crear el nuevo registro
	quote := models.QuotePreview{
		IssuedAt:      fecha,
		Subtotal:      req.Subtotal,
		Tax:           req.Impuesto,
		Total:         req.Total,
		PaymentStatus: models.Pending, // siempre inicia en pending
	}

	// Guardar en la base de datos
	if err := utils.DB.Create(&quote).Error; err != nil {
		log.Printf("ERROR: Fallo al guardar QuotePreview en la base de datos: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar en base de datos"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":            quote.ID,
		"fecha_emision": quote.IssuedAt.Format(time.RFC3339),
		"subtotal":      quote.Subtotal,
		"impuesto":      quote.Tax,
		"total":         quote.Total,
		"estado_pago":   quote.PaymentStatus,
	})
}
