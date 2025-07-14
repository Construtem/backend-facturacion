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

		// Diferenciar entre error de validación y error del sistema
		if strings.Contains(err.Error(), "datos de factura inválidos") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Datos insuficientes para generar factura",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al preparar la factura."})
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
