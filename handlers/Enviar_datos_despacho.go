package handlers

import (
	"backend-facturacion/models"
	"backend-facturacion/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

/*IMPORTANTE :
Para recibir correctamente las peticiones HTTP de esta célula de facturación,
el sistema de despacho debe configurar CORS con los siguientes headers:

CORS Configuration Required:
- Access-Control-Allow-Origin: la URL de producción
- Access-Control-Allow-Methods: POST, OPTIONS
- Access-Control-Allow-Headers: Content-Type, Accept, Authorization
- Access-Control-Allow-Credentials: false

Headers que se envían:
- Content-Type: application/json
- Accept: application/json*/

// DespachoSimpleRequest estructura para envío a despacho
type DespachoSimpleRequest struct {
	CotizacionID      int    `json:"cotizacion_id"`
	RutReceptor       string `json:"rut_receptor"`
	DireccionReceptor string `json:"direccion_receptor"`
	ComunaReceptor    string `json:"comuna_receptor"`
	CiudadReceptor    string `json:"ciudad_receptor"`
}

// EnviarDespachoHandler envía datos a despacho
func EnviarDespachoHandler(c *gin.Context) {
	// Obtener ID desde URL
	quotePreviewIDStr := c.Param("id")
	quotePreviewID, err := strconv.ParseUint(quotePreviewIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de quote preview inválido",
		})
		return
	}

	// Buscar factura en BD
	var factura models.Factura
	if err := utils.DB.Where("quote_preview_id = ?", quotePreviewID).First(&factura).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Factura no encontrada",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error de base de datos",
			})
		}
		return
	}

	// Validar datos mínimos
	if factura.CotizacionID == 0 || factura.DireccionReceptor == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Datos incompletos",
		})
		return
	}

	// Preparar datos para envío
	despachoData := DespachoSimpleRequest{
		CotizacionID:      factura.CotizacionID,
		RutReceptor:       factura.RutReceptor,
		DireccionReceptor: factura.DireccionReceptor,
		ComunaReceptor:    factura.ComunaReceptor,
		CiudadReceptor:    factura.CiudadReceptor,
	}

	// Enviar a despacho
	if err := enviarHTTPADespacho(despachoData, factura.CotizacionID); err != nil {
		log.Printf("Error al enviar a despacho: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al enviar a despacho",
		})
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{
		"message": "Información enviada exitosamente a despacho",
		"data":    despachoData,
	})
}

// enviarHTTPADespacho realiza el envío HTTP
func enviarHTTPADespacho(data DespachoSimpleRequest, cotizacionID int) error {
	baseURL := os.Getenv("BACK_INVENTARIO_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	despachoURL := fmt.Sprintf("%s/api/despachos/%d/ficha", baseURL, cotizacionID)

	// Convertir a JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error al convertir datos a JSON: %w", err)
	}

	// HTTP POST
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("POST", despachoURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error al crear request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error al enviar request: %w", err)
	}
	defer resp.Body.Close()

	// Verificar respuesta
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("despacho respondió con código de error: %d", resp.StatusCode)
	}

	return nil
}
