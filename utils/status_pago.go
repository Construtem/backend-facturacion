package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Request struct {
	status string
}

func MandarStatus(cotizacionID int, status string) error {

	var request1 Request

	jsonData, err := json.Marshal((request1))
	if err != nil {

		return fmt.Errorf("Error al convertir los datos a json: %v", err)

	}

	url := os.Getenv("BACK_VENTAS_URL")

	url1 := fmt.Sprintf("%sapi/status-pago", url)

	req, err := http.NewRequest("POST", url1, bytes.NewBuffer(jsonData))

	// Establecer headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {

		return fmt.Errorf("error al realizar la petición: %v", err)

	}
	defer resp.Body.Close()

	// Verificar el código de respuesta
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {

		return fmt.Errorf("error en la respuesta: código %d", resp.StatusCode)
	}

	fmt.Printf("Status enviado exitosamente: %s\n", status)

	return nil
}
