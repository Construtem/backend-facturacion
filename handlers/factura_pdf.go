package handlers

import (
	services "backend-facturacion/service" // Tu paquete de servicios
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv" // Para parsear el ID

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

	// 1. Obtener o crear los datos de la factura
	factura, err := services.GetOrCreateInvoiceData(uint(quotePreviewID))
	if err != nil {
		if err.Error() == fmt.Sprintf("quote preview con ID %d no encontrada: record not found", quotePreviewID) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cotización preliminar no encontrada."})
			return
		}
		log.Printf("ERROR: Fallo al obtener/crear datos de factura para ID %d: %v", quotePreviewID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al preparar la factura."})
		return
	}

	// 2. Generar el PDF
	var buf bytes.Buffer
	if err := services.GenerateInvoicePDF(factura, &buf); err != nil {
		log.Printf("ERROR: No se pudo generar el PDF: %v", err)
		c.JSON(500, gin.H{"error": "Error interno del servidor"})
		return
	}

	// Verificar que el buffer tenga contenido
	if buf.Len() == 0 {
		log.Printf("ERROR: El PDF generado está vacío")
		c.JSON(500, gin.H{"error": "PDF vacío"})
		return
	}

	// 3. Configurar headers y enviar
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"factura_%s.pdf\"", factura.Folio))
	c.Data(200, "application/pdf", buf.Bytes())
}
