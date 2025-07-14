package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Estructura para el cliente
type Cliente struct {
	Rut      string `json:"rut"`
	Nombre   string `json:"nombre"`
	Telefono string `json:"telefono"`
	Email    string `json:"email"`
}

// Estructura para el usuario
type Usuario struct {
	Nombre string `json:"nombre"`
	Email  string `json:"email"`
	RolID  int    `json:"rol_id"`
}

// Estructura para la dirección
type Direccion struct {
	Direccion string `json:"direccion"`
	Comuna    string `json:"comuna"`
	Ciudad    string `json:"ciudad"`
}

// Estructura para los items
type ItemFactura struct {
	Sku            string  `json:"sku"`
	Nombre         string  `json:"nombre"`
	Cantidad       int     `json:"cantidad"`
	PrecioUnitario float64 `json:"precio_unitario"`
	Subtotal       float64 `json:"subtotal"`
	Sucursal       string  `json:"sucursal"`
}

// Estructura para la respuesta completa
type DatosFactura struct {
	ID           int           `json:"id"`
	FechaCrea    string        `json:"fecha_crea"`
	Estado       string        `json:"estado"`
	CostoEnvio   float64       `json:"costo_envio"`
	TipoDespacho string        `json:"tipo_despacho"`
	Cliente      Cliente       `json:"cliente"`
	Usuario      Usuario       `json:"usuario"`
	Direccion    Direccion     `json:"direccion"`
	Items        []ItemFactura `json:"items"`
	SubtotalNeto float64       `json:"subtotal_neto"`
	Iva          float64       `json:"iva"`
	Total        float64       `json:"total"`
}

func PruebaVentas(c *gin.Context) {
	// Obtener el ID del parámetro de la URL
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Crear datos estáticos para los items
	items := []ItemFactura{
		{
			Sku:            "H001",
			Nombre:         "Martillo carpintero",
			Cantidad:       2,
			PrecioUnitario: 6990,
			Subtotal:       13980,
			Sucursal:       "Bodega Central",
		},
		{
			Sku:            "H002",
			Nombre:         "Destornillador estrella",
			Cantidad:       3,
			PrecioUnitario: 2990,
			Subtotal:       8970,
			Sucursal:       "Bodega Central",
		},
	}

	// Crear la respuesta completa (usando el ID recibido)
	datosFactura := DatosFactura{
		ID:           id, // Usar el ID recibido
		FechaCrea:    "2025-07-01 00:00:00",
		Estado:       "pendiente",
		CostoEnvio:   5000,
		TipoDespacho: "a domicilio",
		Cliente: Cliente{
			Rut:      "11111111-1",
			Nombre:   "Juan Herrera",
			Telefono: "912345111",
			Email:    "juan.herrera@gmail.com",
		},
		Usuario: Usuario{
			Nombre: "Andres Gomez",
			Email:  "agomezr@utem.cl",
			RolID:  2,
		},
		Direccion: Direccion{
			Direccion: "Av. Las Condes 1234",
			Comuna:    "Las Condes",
			Ciudad:    "Santiago",
		},
		Items:        items,
		SubtotalNeto: 22950,
		Iva:          4360.5,
		Total:        32310.5,
	}

	// Devolver los datos en formato JSON
	c.JSON(200, datosFactura)
}
