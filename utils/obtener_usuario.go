package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// ObtenerUsuario obtiene los datos del usuario desde el endpoint
func ObtenerUsuario(cotizacionID int) (*Usuario, error) {
	// URL del endpoint usando variable de entorno
	baseURL := os.Getenv("BACK_VENTAS_URL")
	url := fmt.Sprintf("%sapi/cotizaciones/checkout/%d", baseURL, cotizacionID)

	// Realizar la petición HTTP
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error al llamar al endpoint: %v", err)
	}
	defer resp.Body.Close()

	// Leer la respuesta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer la respuesta: %v", err)
	}

	// Decodificar la respuesta JSON
	var datosFactura DatosFacturaResponse
	if err := json.Unmarshal(body, &datosFactura); err != nil {
		return nil, fmt.Errorf("error al decodificar JSON: %v", err)
	}

	// Retornar únicamente los datos del usuario
	return &datosFactura.Usuario, nil
}
