package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type StatusPagoRequest struct {
	EstadoPago string `json:"estado_pago"`
}

func EstadoPago(status string, CotizacionID int) error {
	baseURL := os.Getenv("BACK_VENTAS_URL")
	url := fmt.Sprintf("%sapi/cotizaciones/status/%d", baseURL, CotizacionID)

	fmt.Printf("🚀 URL: %s\n", url)
	fmt.Printf("📤 Status: %s\n", status)

	requestBody := StatusPagoRequest{
		EstadoPago: status,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("error al convertir datos a JSON: %v", err)
	}

	fmt.Printf("📦 JSON enviado: %s\n", string(jsonData))

	// Prueba con POST primero
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error al crear la petición: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error al realizar la petición: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("error en la respuesta: código %d", resp.StatusCode)
	}

	fmt.Printf("Estado de pago enviado exitosamente: %s para cotización %d\n", status, CotizacionID)
	return nil
}
