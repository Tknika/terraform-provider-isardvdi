package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ListNetworkInterfaces obtiene la lista de todas las interfaces de red
func (c *Client) ListNetworkInterfaces() ([]NetworkInterface, error) {
	reqURL := fmt.Sprintf("https://%s/api/v3/admin/table/interfaces", c.HostURL)
	
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creando petición: %w", err)
	}
	
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando petición: %w", err)
	}

	var interfaces []NetworkInterface
	if err := json.Unmarshal(body, &interfaces); err != nil {
		return nil, fmt.Errorf("error parseando respuesta: %w", err)
	}

	return interfaces, nil
}
