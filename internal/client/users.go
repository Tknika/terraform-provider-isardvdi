package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// User representa un usuario en Isard VDI
type User struct {
	ID                      string                 `json:"id"`
	Name                    string                 `json:"name"`
	Username                string                 `json:"username,omitempty"`
	UID                     string                 `json:"uid,omitempty"`
	Email                   string                 `json:"email,omitempty"`
	Active                  bool                   `json:"active,omitempty"`
	Role                    string                 `json:"role,omitempty"`
	Category                string                 `json:"category,omitempty"`
	Group                   string                 `json:"group,omitempty"`
	SecondaryGroups         []string               `json:"secondary_groups,omitempty"`
	Provider                string                 `json:"provider,omitempty"`
	Accessed                int64                  `json:"accessed,omitempty"`
	EmailVerified           interface{}            `json:"email_verified,omitempty"`
	DisclaimerAcknowledged  interface{}            `json:"disclaimer_acknowledged,omitempty"`
	RoleName                string                 `json:"role_name,omitempty"`
	CategoryName            string                 `json:"category_name,omitempty"`
	GroupName               string                 `json:"group_name,omitempty"`
	SecondaryGroupsNames    []string               `json:"secondary_groups_names,omitempty"`
}

// GetEmailVerified convierte email_verified a bool
func (u *User) GetEmailVerified() bool {
	if u.EmailVerified == nil {
		return false
	}
	switch v := u.EmailVerified.(type) {
	case bool:
		return v
	case float64:
		return v != 0
	case int:
		return v != 0
	default:
		return false
	}
}

// GetDisclaimerAcknowledged convierte disclaimer_acknowledged a bool
func (u *User) GetDisclaimerAcknowledged() bool {
	if u.DisclaimerAcknowledged == nil {
		return false
	}
	switch v := u.DisclaimerAcknowledged.(type) {
	case bool:
		return v
	case float64:
		return v != 0
	case int:
		return v != 0
	default:
		return false
	}
}

// SearchUsersRequest representa la petición de búsqueda de usuarios
type SearchUsersRequest struct {
	Term string `json:"term"`
}

// GetUsers obtiene la lista de usuarios con información completa
func (c *Client) GetUsers() ([]User, error) {
	reqURL := fmt.Sprintf("https://%s/api/v3/admin/users/management/users", c.HostURL)

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
		return nil, fmt.Errorf("error obteniendo usuarios (status %d): %s", res.StatusCode, string(body))
	}

	var users []User
	if err := json.Unmarshal(body, &users); err != nil {
		return nil, fmt.Errorf("error parseando respuesta: %w", err)
	}

	return users, nil
}

// SearchUsers busca usuarios por término (nombre)
func (c *Client) SearchUsers(term string) ([]User, error) {
	reqURL := fmt.Sprintf("https://%s/api/v3/admin/users/search", c.HostURL)

	searchReq := SearchUsersRequest{
		Term: term,
	}

	jsonData, err := json.Marshal(searchReq)
	if err != nil {
		return nil, fmt.Errorf("error creando JSON de búsqueda: %w", err)
	}

	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creando petición POST: %w", err)
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

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error buscando usuarios (status %d): %s", res.StatusCode, string(body))
	}

	var users []User
	if err := json.Unmarshal(body, &users); err != nil {
		return nil, fmt.Errorf("error parseando respuesta: %w", err)
	}

	return users, nil
}

// GetUser obtiene un usuario específico por ID
func (c *Client) GetUser(userID string) (*User, error) {
	reqURL := fmt.Sprintf("https://%s/api/v3/admin/user/%s", c.HostURL, userID)

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
		return nil, fmt.Errorf("error obteniendo usuario (status %d): %s", res.StatusCode, string(body))
	}

	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("error parseando respuesta: %w", err)
	}

	return &user, nil
}
