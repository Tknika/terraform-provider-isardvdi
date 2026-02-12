package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
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
	ISOs       []map[string]interface{} `json:"isos,omitempty"`
	Floppies   []map[string]interface{} `json:"floppies,omitempty"`
}

// CreatePersistentDesktop crea un nuevo persistent desktop
func (c *Client) CreatePersistentDesktop(name, description, templateID string, vcpus *int64, memory *float64, interfaces []string, isos []string, floppies []string) (string, error) {
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
	if vcpus != nil || memory != nil || len(interfaces) > 0 || len(isos) > 0 || len(floppies) > 0 {
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
		if len(isos) > 0 {
			isoList := make([]map[string]interface{}, len(isos))
			for i, isoID := range isos {
				isoList[i] = map[string]interface{}{"id": isoID}
			}
			hardware["isos"] = isoList
		}
		if len(floppies) > 0 {
			floppyList := make([]map[string]interface{}, len(floppies))
			for i, floppyID := range floppies {
				floppyList[i] = map[string]interface{}{"id": floppyID}
			}
			hardware["floppies"] = floppyList
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

// StopDesktop detiene un desktop
func (c *Client) StopDesktop(desktopID string) error {
	reqURL := fmt.Sprintf("https://%s/api/v3/desktop/stop/%s", c.HostURL, desktopID)

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return fmt.Errorf("error creando la petición GET: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error ejecutando GET: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error leyendo respuesta: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error deteniendo desktop (status %d): %s", res.StatusCode, string(body))
	}

	return nil
}

// GetDesktopStatus obtiene el estado actual de un desktop
func (c *Client) GetDesktopStatus(desktopID string) (string, error) {
	reqURL := fmt.Sprintf("https://%s/api/v3/domain/info/%s", c.HostURL, desktopID)

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("error creando la petición GET: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error ejecutando GET: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error leyendo respuesta: %w", err)
	}

	if res.StatusCode == http.StatusNotFound {
		return "not_found", nil
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error obteniendo estado del desktop (status %d): %s", res.StatusCode, string(body))
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("error parseando respuesta JSON: %w", err)
	}

	if status, ok := response["status"].(string); ok {
		return status, nil
	}

	return "unknown", nil
}

// WaitForDesktopStopped espera a que un desktop se detenga completamente
func (c *Client) WaitForDesktopStopped(desktopID string, maxWaitSeconds int) error {
	// Verificar inmediatamente si ya está detenido (antes de esperar)
	status, err := c.GetDesktopStatus(desktopID)
	if err != nil {
		return fmt.Errorf("error obteniendo estado del desktop: %w", err)
	}
	
	// Estados que indican que el desktop está detenido
	if status == "Stopped" || status == "stopped" || status == "Shutdown" || status == "shutdown" || status == "Failed" || status == "failed" {
		return nil
	}
	
	// Si no está detenido, esperar con polling
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	
	timeout := time.After(time.Duration(maxWaitSeconds) * time.Second)
	
	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout esperando a que se detenga el desktop después de %d segundos", maxWaitSeconds)
		case <-ticker.C:
			status, err := c.GetDesktopStatus(desktopID)
			if err != nil {
				return fmt.Errorf("error obteniendo estado del desktop: %w", err)
			}
			
			// Estados que indican que el desktop está detenido
			if status == "Stopped" || status == "stopped" || status == "Shutdown" || status == "shutdown" || status == "Failed" || status == "failed" {
				return nil
			}
		}
	}
}

// ForceStopDesktop fuerza la parada de un desktop usando el endpoint admin
func (c *Client) ForceStopDesktop(desktopID string) error {
	reqURL := fmt.Sprintf("https://%s/api/v3/admin/multiple_actions", c.HostURL)

	payload := map[string]interface{}{
		"ids":    []string{desktopID},
		"action": "stopping",
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

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error forzando parada del desktop (status %d): %s", res.StatusCode, string(body))
	}

	return nil
}
