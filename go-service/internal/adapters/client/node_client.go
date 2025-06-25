package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/luke385/skill-test/internal/domain"
)

type NodeAPIClient struct {
	baseURL         string
	httpClient      *http.Client
	csrfCookieName  string
	csrfCookieValue string
	accessToken     string
	refreshToken    string
	cookieHeader    string
}

// NewNodeAPIClient authenticates against the NodeJS backend and prepares an authorized client.
func NewNodeAPIClient() (*NodeAPIClient, error) {
	baseURL := os.Getenv("NODE_API_URL")
	user := os.Getenv("NODE_API_USER")
	pass := os.Getenv("NODE_API_PASS")
	csrf := os.Getenv("CSRF_COOKIE_NAME")
	if baseURL == "" || user == "" || pass == "" || csrf == "" {
		return nil, fmt.Errorf("missing required environment variables")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	api := &NodeAPIClient{
		baseURL:        baseURL,
		httpClient:     client,
		csrfCookieName: csrf,
	}

	// Authenticate with backend
	loginURL := fmt.Sprintf("%s/api/v1/auth/login", baseURL)
	payload := map[string]string{"username": user, "password": pass}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(payload); err != nil {
		return nil, fmt.Errorf("failed to encode login payload: %w", err)
	}

	req, err := http.NewRequest("POST", loginURL, buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	log.Println("Attempting login against NodeJS backend...")
	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()
	log.Printf("Login status: %d", resp.StatusCode)

	// Parse cookies and extract tokens
	var accessToken, refreshToken, csrfToken string
	for _, ck := range resp.Cookies() {
		log.Printf("Set-Cookie: %s=%s", ck.Name, ck.Value)
		switch ck.Name {
		case "accessToken":
			accessToken = ck.Value
		case "refreshToken":
			refreshToken = ck.Value
		case csrf:
			csrfToken = ck.Value
		}
	}

	// Build a single Cookie header
	api.cookieHeader = strings.Join([]string{
		fmt.Sprintf("accessToken=%s", accessToken),
		fmt.Sprintf("refreshToken=%s", refreshToken),
		fmt.Sprintf("%s=%s", csrf, csrfToken),
	}, "; ")
	api.accessToken = accessToken
	api.refreshToken = refreshToken
	api.csrfCookieValue = csrfToken

	// error handling
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("login failed %d: %s", resp.StatusCode, string(b))
	}
	if api.accessToken == "" {
		return nil, fmt.Errorf("accessToken cookie not found in login response")
	}
	if api.csrfCookieValue == "" {
		return nil, fmt.Errorf("CSRF cookie %q not found in login response", api.csrfCookieName)
	}
	return api, nil
}

// GetByID fetches a student by ID, using authentication headers and cookies as required by backend
func (c *NodeAPIClient) GetByID(id string) (*domain.Student, error) {
	fullURL := fmt.Sprintf("%s/api/v1/students/%s", c.baseURL, id)
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("X-CSRF-Token", c.csrfCookieValue)
	req.Header.Set("Cookie", c.cookieHeader)

	// Debug logging
	log.Printf("   GET %s", fullURL)
	log.Printf("   Authorization: Bearer %s", c.accessToken)
	log.Printf("   X-CSRF-Token: %s", c.csrfCookieValue)
	log.Printf("   Cookie: %s", c.cookieHeader)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("backend request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read backend response: %w", err)
	}

	log.Printf("STATUS: %d", resp.StatusCode)
	log.Printf("BODY: %s", string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("backend %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var student domain.Student
	if err := json.Unmarshal(bodyBytes, &student); err != nil {
		return nil, fmt.Errorf("failed to decode backend JSON: %w", err)
	}
	return &student, nil
}
