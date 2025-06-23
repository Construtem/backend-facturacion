package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func VerificarPago(id int) string {
	accessToken := os.Getenv("ACCESS_TOKEN")
	url := fmt.Sprintf("https://api.mercadopago.com/v1/payments/%d", id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "error"
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "error"
	}
	defer resp.Body.Close()

	var result struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "error"
	}

	return result.Status
}
