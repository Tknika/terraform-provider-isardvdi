package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// NetworkInterface representa la estructura de una interfaz de red en la API
type NetworkInterface struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Net         string                 `json:"net"`
	Kind        string                 `json:"kind,omitempty"`
	Model       string                 `json:"model,omitempty"`
	QoSID       string                 `json:"qos_id,omitempty"`
	Ifname      string                 `json:"ifname,omitempty"`
	Allowed     map[string]interface{} `json:"allowed,omitempty"`
}

// CreateNetworkInterface crea una nueva interfaz de red
func (c *Client) CreateNetworkInterface(id, name, description, net, kind, model, qosID, ifname string, allowed map[string]interface{}) error {
	reqURL := fmt.Sprintf("https://%s/api/v3/admin/table/add/interfaces", c.HostURL)

	// Construir el payload
	payload := map[string]interface{}{
		"id":   id,
		"name": name,
		"net":  net,
	}
	
	if description != "" {
		payload["description"] = description
	}
	
	if kind != "" {
		payload["kind"] = kind
	}
	
	if model != "" {
		payload["model"] = model
	}
	
	if qosID != "" {
		payload["qos_id"] = qosID
	}
	
	if ifname != "" {
		payload["ifname"] = ifname
	}
	
	if allowed != nil {
		payload["allowed"] = allowed
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error codificando JSON: %w", err)
	}

	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creando la petición POST: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error ejecutando POST: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error leyendo respuesta: %w", err)
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		return fmt.Errorf("error creando interfaz de red (status %d): %s", res.StatusCode, string(body))
	}

	return nil
}

// GetNetworkInterface obtiene la información de una interfaz de red
func (c *Client) GetNetworkInterface(interfaceID string) (*NetworkInterface, error) {
	reqURL := fmt.Sprintf("https://%s/api/v3/admin/table/interfaces", c.HostURL)

	// Crear payload con el ID para obtener un item específico
	payload := map[string]interface{}{
		"id": interfaceID,
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
		return nil, fmt.Errorf("network interface not found")
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error obteniendo interfaz de red (status %d): %s", res.StatusCode, string(body))
	}

	// Parsear la respuesta
	var iface NetworkInterface
	if err := json.Unmarshal(body, &iface); err != nil {
		return nil, fmt.Errorf("error parseando respuesta JSON: %w", err)
	}

	return &iface, nil
}

// UpdateNetworkInterface actualiza una interfaz de red existente
func (c *Client) UpdateNetworkInterface(id string, name, description, net, kind, model, qosID, ifname *string, allowed map[string]interface{}) error {
	reqURL := fmt.Sprintf("https://%s/api/v3/admin/table/update/interfaces", c.HostURL)

	// Construir el payload con el ID y los campos a actualizar
	payload := map[string]interface{}{
		"id": id,
	}
	
	if name != nil {
		payload["name"] = *name
	}
	
	if description != nil {
		payload["description"] = *description
	}
	
	if net != nil {
		payload["net"] = *net
	}
	
	if kind != nil {
		payload["kind"] = *kind
	}
	
	if model != nil {
		payload["model"] = *model
	}
	
	if qosID != nil {
		payload["qos_id"] = *qosID
	}
	
	if ifname != nil {
		payload["ifname"] = *ifname
	}
	
	if allowed != nil {
		payload["allowed"] = allowed
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
		return fmt.Errorf("error actualizando interfaz de red (status %d): %s", res.StatusCode, string(body))
	}

	return nil
}

// DeleteNetworkInterface elimina una interfaz de red
func (c *Client) DeleteNetworkInterface(interfaceID string) error {
	reqURL := fmt.Sprintf("https://%s/api/v3/admin/table/interfaces/%s", c.HostURL, interfaceID)

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

	return fmt.Errorf("error eliminando interfaz de red (status %d): %s", res.StatusCode, string(body))
}
