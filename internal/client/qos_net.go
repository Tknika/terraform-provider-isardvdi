package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// QoSNet representa la estructura de un QoS de red en la API
type QoSNet struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Bandwidth   map[string]interface{} `json:"bandwidth,omitempty"`
}

// CreateQoSNet crea un nuevo QoS de red
func (c *Client) CreateQoSNet(name, description string, bandwidth map[string]interface{}) (string, error) {
	reqURL := fmt.Sprintf("https://%s/api/v3/admin/table/add/qos_net", c.HostURL)

	// Construir el payload
	payload := map[string]interface{}{
		"name": name,
	}
	
	if description != "" {
		payload["description"] = description
	}
	
	if bandwidth != nil {
		payload["bandwidth"] = bandwidth
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error codificando JSON: %w", err)
	}

	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creando la petición POST: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error ejecutando POST: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error leyendo respuesta: %w", err)
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("error creando QoS de red (status %d): %s", res.StatusCode, string(body))
	}

	// La API devuelve el ID en el campo 'id' o podemos usar el nombre como ID
	// Primero intentamos obtener el ID de la base de datos
	qos, err := c.GetQoSNet(name)
	if err != nil {
		// Si no podemos obtenerlo, usamos el nombre como ID
		return name, nil
	}

	return qos.ID, nil
}

// GetQoSNet obtiene la información de un QoS de red
func (c *Client) GetQoSNet(qosID string) (*QoSNet, error) {
	reqURL := fmt.Sprintf("https://%s/api/v3/admin/table/qos_net", c.HostURL)

	// Crear payload con el ID para obtener un item específico
	payload := map[string]interface{}{
		"id": qosID,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error codificando JSON: %w", err)
	}

	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creando la petición POST: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando POST: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta: %w", err)
	}

	if res.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("qos_net not found")
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error obteniendo QoS de red (status %d): %s", res.StatusCode, string(body))
	}

	// Parsear la respuesta
	var qos QoSNet
	if err := json.Unmarshal(body, &qos); err != nil {
		return nil, fmt.Errorf("error parseando respuesta JSON: %w", err)
	}

	return &qos, nil
}

// UpdateQoSNet actualiza un QoS de red existente
func (c *Client) UpdateQoSNet(qosID string, name, description *string, bandwidth map[string]interface{}) error {
	reqURL := fmt.Sprintf("https://%s/api/v3/admin/table/update/qos_net", c.HostURL)

	// Construir el payload con el ID y los campos a actualizar
	payload := map[string]interface{}{
		"id": qosID,
	}
	
	if name != nil {
		payload["name"] = *name
	}
	
	if description != nil {
		payload["description"] = *description
	}
	
	if bandwidth != nil {
		payload["bandwidth"] = bandwidth
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error codificando JSON: %w", err)
	}

	req, err := http.NewRequest("PUT", reqURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creando la petición PUT: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error ejecutando PUT: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error leyendo respuesta: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error actualizando QoS de red (status %d): %s", res.StatusCode, string(body))
	}

	return nil
}

// DeleteQoSNet elimina un QoS de red
func (c *Client) DeleteQoSNet(qosID string) error {
	reqURL := fmt.Sprintf("https://%s/api/v3/admin/table/qos_net/%s", c.HostURL, qosID)

	req, err := http.NewRequest("DELETE", reqURL, nil)
	if err != nil {
		return fmt.Errorf("error creando la petición DELETE: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error ejecutando DELETE: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error leyendo respuesta: %w", err)
	}

	// Considerar éxito los códigos 200, 204 (No Content) y 404 (ya no existe)
	if res.StatusCode == http.StatusOK || res.StatusCode == http.StatusNoContent || res.StatusCode == http.StatusNotFound {
		return nil
	}

	return fmt.Errorf("error eliminando QoS de red (status %d): %s", res.StatusCode, string(body))
}
