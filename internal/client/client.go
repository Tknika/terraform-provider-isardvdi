package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/tknika/terraform-provider-isardvdi/internal/constants"
)

// Client holds the connection information
type Client struct {
	HTTPClient *http.Client
	HostURL    string
	Token      string
}

// NewClient creates a new client
func NewClient(host, token string, sslVerification bool) *Client {
	// Configurar transporte HTTP con opción de verificación SSL configurable
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !sslVerification},
	}
	
	return &Client{
		HTTPClient: &http.Client{
			Timeout:   60 * time.Second,
			Transport: tr,
		},
		HostURL: host,
		Token:   token,
	}
}

// SignIn performs the authentication flow
func (c *Client) SignIn(authMethod, categoryID, username, password string) error {
	if authMethod == "token" {
		// Cuando usamos token, simplemente lo usamos directamente sin hacer llamadas adicionales
		// El token ya está almacenado en c.Token desde NewClient
		return nil
	}

	if authMethod == "form" {
		// Construir URL con query params
		reqURL := fmt.Sprintf("https://%s%s?provider=form&category_id=%s", c.HostURL, constants.LoginPath, categoryID)

		// Multipart form data
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("username", username)
		_ = writer.WriteField("password", password)
		err := writer.Close()
		if err != nil {
			return err
		}

		req, err := http.NewRequest("POST", reqURL, body)
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Accept", "text/plain")

		println(req.URL.String())

		return c.executeAuthRequest(req)
	}

	return nil
}

func (c *Client) executeAuthRequest(req *http.Request) error {
	body, err := c.doRequest(req)
	if err != nil {
		return err
	}

	// Intentamos parsear la respuesta para encontrar el token temporal (JSON)
	var authResp map[string]interface{}
	if err := json.Unmarshal(body, &authResp); err == nil {
		// Búsqueda del token en la respuesta
		// Caso 1: {"data": "token_string"}
		if token, ok := authResp["data"].(string); ok {
			c.Token = token
			return nil
		}
		// Caso 2: {"token": "token_string"}
		if token, ok := authResp["token"].(string); ok {
			c.Token = token
			return nil
		}

		// Caso 3: {"data": {"token": "token_string"}}
		if data, ok := authResp["data"].(map[string]interface{}); ok {
			if token, ok := data["token"].(string); ok {
				c.Token = token
				return nil
			}
		}
	}

	// Si falla el parseo JSON o no se encuentra estructura,
	// y el body no está vacío, asumimos que es el token en texto plano (auth form)
	if len(body) > 0 {
		// Podríamos añadir validación extra (ej. longitud mínima)
		c.Token = string(body)
		return nil
	}

	return fmt.Errorf("no se encontró el token en la respuesta de login. Respuesta cruda: %s", string(body))
}

// doRequest helper for executing requests
func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, nil
}
