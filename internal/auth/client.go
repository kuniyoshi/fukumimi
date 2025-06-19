package auth

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"syscall"

	"github.com/kuniyoshi/fukumimi/internal/config"
	"golang.org/x/term"
)

const (
	// UserAgent identifies this client
	UserAgent = "Fukumimi/0.1.0 (https://github.com/kuniyoshi/fukumimi)"
)

type Client struct {
	httpClient *http.Client
	jar        *cookiejar.Jar
}

type storedCookie struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Domain   string `json:"domain"`
	Path     string `json:"path"`
	Secure   bool   `json:"secure"`
	HttpOnly bool   `json:"httpOnly"`
}

// customTransport adds User-Agent header to all requests
type customTransport struct {
	base http.RoundTripper
}

func (t *customTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Set User-Agent header
	req.Header.Set("User-Agent", UserAgent)

	return t.base.RoundTrip(req)
}

func NewClient() *Client {
	jar, _ := cookiejar.New(nil)

	// Create HTTP client with custom transport to set User-Agent
	transport := &customTransport{
		base: http.DefaultTransport,
	}
	client := &http.Client{
		Jar:       jar,
		Transport: transport,
	}

	return &Client{
		httpClient: client,
		jar:        jar,
	}
}

func (c *Client) Login() error {
	// Try to load existing cookies first
	if err := c.loadCookies(); err == nil {
		// Test if cookies are still valid
		if c.isAuthenticated() {
			return nil
		}
	}

	// Get credentials from user
	username, password, err := c.getCredentials()
	if err != nil {
		return err
	}

	// Perform login
	formData := url.Values{
		"username": {username},
		"password": {password},
	}

	resp, err := c.httpClient.PostForm(config.LoginURL, formData)
	if err != nil {
		return fmt.Errorf("failed to submit login form: %w", err)
	}
	defer resp.Body.Close()

	// Check if login was successful
	if !c.isAuthenticated() {
		return fmt.Errorf("authentication failed - please check your credentials")
	}

	// Save cookies for future use
	if err := c.saveCookies(); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	return nil
}

func (c *Client) Get(url string) (*http.Response, error) {
	return c.httpClient.Get(url)
}

func (c *Client) getCredentials() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}
	username = strings.TrimSpace(username)

	fmt.Print("Password: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}
	fmt.Println() // New line after password input
	password := string(passwordBytes)

	return username, password, nil
}

func (c *Client) isAuthenticated() bool {
	// Make a test request to check if we're authenticated
	// This is a simplified check - adjust based on actual site behavior
	u, _ := url.Parse(config.LoginURL)
	cookies := c.jar.Cookies(u)

	// Check if we have session cookies
	for _, cookie := range cookies {
		if strings.Contains(strings.ToLower(cookie.Name), "session") ||
			strings.Contains(strings.ToLower(cookie.Name), "auth") {
			return true
		}
	}

	return false
}

func (c *Client) saveCookies() error {
	u, _ := url.Parse(config.LoginURL)
	cookies := c.jar.Cookies(u)

	var storedCookies []storedCookie
	for _, cookie := range cookies {
		storedCookies = append(storedCookies, storedCookie{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Domain:   cookie.Domain,
			Path:     cookie.Path,
			Secure:   cookie.Secure,
			HttpOnly: cookie.HttpOnly,
		})
	}

	data, err := json.MarshalIndent(storedCookies, "", "  ")
	if err != nil {
		return err
	}

	cookiePath := config.GetCookieFilePath()
	return os.WriteFile(cookiePath, data, 0600)
}

func (c *Client) loadCookies() error {
	cookiePath := config.GetCookieFilePath()
	data, err := os.ReadFile(cookiePath)
	if err != nil {
		return err
	}

	var storedCookies []storedCookie
	if err := json.Unmarshal(data, &storedCookies); err != nil {
		return err
	}

	u, _ := url.Parse(config.LoginURL)
	var cookies []*http.Cookie
	for _, sc := range storedCookies {
		cookies = append(cookies, &http.Cookie{
			Name:     sc.Name,
			Value:    sc.Value,
			Domain:   sc.Domain,
			Path:     sc.Path,
			Secure:   sc.Secure,
			HttpOnly: sc.HttpOnly,
		})
	}

	c.jar.SetCookies(u, cookies)
	return nil
}
