package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Template representa la estructura de un template en la API
type Template struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Group       string `json:"group"`
	UserID      string `json:"user_id"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	Status      string `json:"status"`
	DesktopSize int64  `json:"desktop_size"`
}

// GetTemplates obtiene la lista de templates disponibles para el usuario
func (c *Client) GetTemplates() ([]Template, error) {
	reqURL := fmt.Sprintf("https://%s/api/v3/user/templates", c.HostURL)

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var templates []Template
	if err := json.Unmarshal(body, &templates); err != nil {
		return nil, err
	}

	return templates, nil
}
