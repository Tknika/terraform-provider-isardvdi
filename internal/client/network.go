package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// Network representa la estructura de una red en la API
type Network struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Model       string                 `json:"model"`
	QoSID       string                 `json:"qos_id"`
	MetadataID  string                 `json:"-"` // No se parsea directamente del JSON
	Allowed     map[string]interface{} `json:"allowed"`
	User        string                 `json:"user"`
	Group       string                 `json:"group"`
	Category    string                 `json:"category"`
	Created     string                 `json:"created"`
	Modified    string                 `json:"modified"`
}

// CreateNetwork crea una nueva red de usuario
func (c *Client) CreateNetwork(name, description, model, qosID string, allowed map[string]interface{}) (string, error) {
	reqURL := fmt.Sprintf("https://%s/api/v3/user/networks", c.HostURL)

	// Construir el payload
	payload := map[string]interface{}{
		"name":        name,
		"description": description,
	}
	
	if model != "" {
		payload["model"] = model
	}
	
	if qosID != "" {
		payload["qos_id"] = qosID
	}
	
	if allowed != nil {
		payload["allowed"] = allowed
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
		return "", fmt.Errorf("error creando red (status %d): %s", res.StatusCode, string(body))
	}

	// Parsear la respuesta para obtener el ID
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("error parseando respuesta JSON: %w", err)
	}

	networkID, ok := response["id"].(string)
	if !ok {
		return "", fmt.Errorf("no se encontró el ID en la respuesta: %s", string(body))
	}

	return networkID, nil
}

// GetNetwork obtiene la información de una red
func (c *Client) GetNetwork(networkID string) (*Network, error) {
	reqURL := fmt.Sprintf("https://%s/api/v3/user/networks/%s", c.HostURL, networkID)

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creando la petición GET: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando GET: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta: %w", err)
	}

	if res.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("network not found")
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error obteniendo red (status %d): %s", res.StatusCode, string(body))
	}

	// Parsear la respuesta usando un decoder con UseNumber para manejar números grandes
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()
	
	var rawNetwork map[string]interface{}
	if err := decoder.Decode(&rawNetwork); err != nil {
		return nil, fmt.Errorf("error parseando respuesta JSON: %w", err)
	}

	network := &Network{}
	
	// Parsear campos uno por uno
	if id, ok := rawNetwork["id"].(string); ok {
		network.ID = id
	}
	if name, ok := rawNetwork["name"].(string); ok {
		network.Name = name
	}
	if desc, ok := rawNetwork["description"].(string); ok {
		network.Description = desc
	}
	if model, ok := rawNetwork["model"].(string); ok {
		network.Model = model
	}
	if qosID, ok := rawNetwork["qos_id"].(string); ok {
		network.QoSID = qosID
	}
	
	// Parsear metadata_id como json.Number para manejar valores grandes
	// Convertir a string sin notación científica
	if metadataIDNum, ok := rawNetwork["metadata_id"].(json.Number); ok {
		// Intentar parsear como uint64 desde el string del json.Number
		if val, err := strconv.ParseUint(metadataIDNum.String(), 10, 64); err == nil {
			network.MetadataID = fmt.Sprintf("%d", val)
		} else {
			// Si falla (notación científica), intentar como float y convertir
			if fval, err := metadataIDNum.Float64(); err == nil {
				network.MetadataID = fmt.Sprintf("%.0f", fval)
			} else {
				network.MetadataID = metadataIDNum.String()
			}
		}
	}

	return network, nil
}

// UpdateNetwork actualiza una red existente
func (c *Client) UpdateNetwork(networkID string, name, description, qosID *string, allowed map[string]interface{}) error {
	reqURL := fmt.Sprintf("https://%s/api/v3/user/networks/%s", c.HostURL, networkID)

	// Construir el payload solo con los campos que se actualizan
	payload := make(map[string]interface{})
	
	if name != nil {
		payload["name"] = *name
	}
	
	if description != nil {
		payload["description"] = *description
	}
	
	if qosID != nil {
		payload["qos_id"] = *qosID
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
		return fmt.Errorf("error actualizando red (status %d): %s", res.StatusCode, string(body))
	}

	return nil
}

// DeleteNetwork elimina una red
func (c *Client) DeleteNetwork(networkID string) error {
	reqURL := fmt.Sprintf("https://%s/api/v3/user/networks/%s", c.HostURL, networkID)

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

	return fmt.Errorf("error eliminando red (status %d): %s", res.StatusCode, string(body))
}
