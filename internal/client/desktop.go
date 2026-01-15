package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Desktop representa la estructura de un desktop en la API
type Desktop struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	TemplateID  string  `json:"template_id"`
	VCPUs       int64   `json:"vcpus,omitempty"`
	Memory      float64 `json:"memory,omitempty"`
}

// HardwareSpec especifica el hardware personalizado para un desktop
type HardwareSpec struct {
	VCPUs      *int64   `json:"vcpus,omitempty"`
	Memory     *float64 `json:"memory,omitempty"`
	DiskBus    string   `json:"disk_bus,omitempty"`
	BootOrder  []string `json:"boot_order,omitempty"`
	Graphics   []string `json:"graphics,omitempty"`
	Videos     []string `json:"videos,omitempty"`
	Interfaces []string `json:"interfaces,omitempty"`
}

// CreatePersistentDesktop crea un nuevo persistent desktop
func (c *Client) CreatePersistentDesktop(name, description, templateID string, vcpus *int64, memory *float64, interfaces []string) (string, error) {
	reqURL := fmt.Sprintf("https://%s/api/v3/persistent_desktop", c.HostURL)

	// Construir el payload
	payload := map[string]interface{}{
		"name":        name,
		"template_id": templateID,
	}
	
	if description != "" {
		payload["description"] = description
	}

	// Agregar hardware personalizado si se especifica
	if vcpus != nil || memory != nil || len(interfaces) > 0 {
		hardware := make(map[string]interface{})
		if vcpus != nil {
			hardware["vcpus"] = *vcpus
		}
		if memory != nil {
			hardware["memory"] = *memory
		}
		if len(interfaces) > 0 {
			hardware["interfaces"] = interfaces
		}
		payload["hardware"] = hardware
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
		return "", fmt.Errorf("error creando desktop (status %d): %s", res.StatusCode, string(body))
	}

	// Parsear la respuesta para obtener el ID
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("error parseando respuesta JSON: %w", err)
	}

	desktopID, ok := response["id"].(string)
	if !ok {
		return "", fmt.Errorf("no se encontró el ID en la respuesta: %s", string(body))
	}

	return desktopID, nil
}

// GetDesktop obtiene la información de un desktop
func (c *Client) GetDesktop(desktopID string) (*Desktop, error) {
	reqURL := fmt.Sprintf("https://%s/api/v3/domain/info/%s", c.HostURL, desktopID)

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
		return nil, fmt.Errorf("desktop not found")
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error obteniendo desktop (status %d): %s", res.StatusCode, string(body))
	}

	// Parsear la respuesta
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error parseando respuesta JSON: %w", err)
	}

	desktop := &Desktop{
		ID: desktopID,
	}

	if name, ok := response["name"].(string); ok {
		desktop.Name = name
	}
	if desc, ok := response["description"].(string); ok {
		desktop.Description = desc
	}
	if createDict, ok := response["create_dict"].(map[string]interface{}); ok {
		if origin, ok := createDict["origin"].(string); ok {
			desktop.TemplateID = origin
		}
	}
	
	// Leer el hardware
	if hardware, ok := response["hardware"].(map[string]interface{}); ok {
		if vcpus, ok := hardware["vcpus"].(float64); ok {
			desktop.VCPUs = int64(vcpus)
		}
		if memory, ok := hardware["memory"].(float64); ok {
			desktop.Memory = memory
		}
	}

	return desktop, nil
}

// DeleteDesktop deletes a desktop by its ID
func (c *Client) DeleteDesktop(desktopID string) error {
	reqURL := fmt.Sprintf("https://%s/api/v3/desktop/%s/true", c.HostURL, desktopID)

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

	// Leer el body para obtener información de error si es necesario
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error leyendo respuesta: %w", err)
	}

	// Considerar éxito los códigos 200, 204 (No Content) y 404 (ya no existe)
	if res.StatusCode == http.StatusOK || res.StatusCode == http.StatusNoContent || res.StatusCode == http.StatusNotFound {
		return nil
	}

	return fmt.Errorf("error eliminando desktop (status %d): %s", res.StatusCode, string(body))
}
