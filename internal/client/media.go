package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Media representa un media en Isard VDI
type Media struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	URL         string                 `json:"url-web,omitempty"`
	URLIsard    interface{}            `json:"url-isard,omitempty"`
	Kind        string                 `json:"kind"`
	Status      string                 `json:"status,omitempty"`
	User        string                 `json:"user"`
	Category    string                 `json:"category"`
	Group       string                 `json:"group"`
	Allowed     map[string]interface{} `json:"allowed,omitempty"`
	Icon        string                 `json:"icon,omitempty"`
	Path        string                 `json:"path,omitempty"`
	Progress    map[string]interface{} `json:"progress,omitempty"`
	Accessed    float64                `json:"accessed,omitempty"`
}

// CreateMedia crea un nuevo media
func (c *Client) CreateMedia(
	name string,
	description string,
	url string,
	kind string,
	allowed map[string]interface{},
) (string, error) {
	reqURL := fmt.Sprintf("https://%s/api/v3/media", c.HostURL)

	payload := map[string]interface{}{
		"name":        name,
		"description": description,
		"url":         url,
		"kind":        kind,
	}

	// Añadir allowed si se especifica
	if allowed != nil && len(allowed) > 0 {
		payload["allowed"] = allowed
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error creando JSON: %w", err)
	}

	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creando petición POST: %w", err)
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

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error creando media (status %d): %s", res.StatusCode, string(body))
	}

	// La API devuelve el media creado, extraer el ID
	// Como la API puede devolver vacío, necesitamos obtener el media por nombre
	// Esperamos un poco para que se cree
	time.Sleep(2 * time.Second)
	
	// Obtener la lista de medias y buscar el que acabamos de crear
	medias, err := c.GetMedias()
	if err != nil {
		return "", fmt.Errorf("error obteniendo medias después de crear: %w", err)
	}

	// Buscar el media recién creado por nombre
	for _, media := range medias {
		if media.Name == name {
			return media.ID, nil
		}
	}

	return "", fmt.Errorf("no se pudo encontrar el media creado con nombre: %s", name)
}

// GetMedia obtiene información de un media específico
func (c *Client) GetMedia(mediaID string) (*Media, error) {
	// Obtener todos los medias y buscar el específico
	medias, err := c.GetMedias()
	if err != nil {
		return nil, err
	}

	for _, media := range medias {
		if media.ID == mediaID {
			return &media, nil
		}
	}

	return nil, fmt.Errorf("media no encontrado: %s", mediaID)
}

// GetMedias obtiene la lista de medias del usuario
func (c *Client) GetMedias() ([]Media, error) {
	reqURL := fmt.Sprintf("https://%s/api/v3/media", c.HostURL)

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creando petición GET: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando GET: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error obteniendo medias (status %d): %s", res.StatusCode, string(body))
	}

	var medias []Media
	if err := json.Unmarshal(body, &medias); err != nil {
		return nil, fmt.Errorf("error parseando respuesta: %w", err)
	}

	return medias, nil
}

// DeleteMedia elimina un media
func (c *Client) DeleteMedia(mediaID string) error {
	reqURL := fmt.Sprintf("https://%s/api/v3/media/%s", c.HostURL, mediaID)

	req, err := http.NewRequest("DELETE", reqURL, nil)
	if err != nil {
		return fmt.Errorf("error creando petición DELETE: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error ejecutando DELETE: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error leyendo respuesta: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error eliminando media (status %d): %s", res.StatusCode, string(body))
	}

	return nil
}
