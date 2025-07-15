package handlers

import (
	services "backend-facturacion/service" // Tu paquete de servicios
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv" // Para parsear el ID
	"strings"

	"github.com/gin-gonic/gin"
)

// GenerateInvoicePDFHandler genera y sirve el PDF de la factura
func GenerateInvoicePDFHandler(c *gin.Context) {
	// Obtener el ID de la preview de cotización desde la URL
	quotePreviewIDStr := c.Param("id")
	quotePreviewID, err := strconv.ParseUint(quotePreviewIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cotización inválido."})
		return
	}

	// Obtener datos de factura
	factura, err := services.GetOrCreateInvoiceData(uint(quotePreviewID))
	if err != nil {
		log.Printf("ERROR: %v", err)

		// Diferenciar tipos de error para respuestas más específicas
		if strings.Contains(err.Error(), "factura no encontrada") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Factura no encontrada",
				"details": fmt.Sprintf("No existe una factura para quote_preview_id %d", quotePreviewID),
				"action":  "Debe crear la factura primero desde el sistema de cotizaciones",
			})
			return
		}

		if strings.Contains(err.Error(), "datos de factura incompletos") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Datos de factura incompletos",
				"details": err.Error(),
				"action":  "Complete los datos faltantes antes de generar el PDF",
			})
			return
		}

		// Error genérico del sistema
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	// Generar el PDF
	var buf bytes.Buffer
	if err := services.GenerateInvoicePDF(factura, &buf); err != nil {
		log.Printf("ERROR: No se pudo generar el PDF: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	// Verificar que el buffer tenga contenido
	if buf.Len() == 0 {
		log.Printf("ERROR: El PDF generado está vacío")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "PDF vacío"})
		return
	}

	// Configurar headers y enviar
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"factura_%s.pdf\"", factura.Folio))
	c.Data(http.StatusOK, "application/pdf", buf.Bytes())
}
