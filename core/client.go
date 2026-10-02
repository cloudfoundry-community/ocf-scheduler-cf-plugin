package core

import (
	"bytes"
	"crypto/tls"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ess/hype"
)

type Driver struct {
	baseURL     string
	raw         *hype.Driver
	token       string
	accept      *hype.Header
	contentType *hype.Header
	auth        *hype.Header
	userAgent   *hype.Header
}

var Client *Driver

// httpClient matches hype's 20 s timeout, so a hung scheduler fails
// instead of hanging the command.
var httpClient = &http.Client{Timeout: 20 * time.Second}

func NewDriver(baseURL string, token string) (*Driver, error) {
	// TODO: figure out how to make this configurable
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	raw, err := hype.New(baseURL)
	if err != nil {
		return nil, err
	}

	d := &Driver{
		baseURL,
		raw,
		token,
		hype.NewHeader("Accept", "application/json"),
		hype.NewHeader("Content-Type", "application/json"),
		hype.NewHeader("Authorization", token),
		hype.NewHeader("User-Agent", "ocf-scheduler-cf-plugin"),
	}

	return d, nil
}

func (driver *Driver) Token() string {
	return driver.token
}

func (driver *Driver) Delete(path string, params hype.Params) hype.Response {
	return driver.
		raw.
		Delete(path, params).
		WithHeaderSet(driver.accept, driver.contentType, driver.auth).
		Response()
}

func (driver *Driver) Get(path string, params hype.Params) hype.Response {
	raw := driver.raw
	get := raw.Get(path, params)
	withHeaders := get.WithHeaderSet(driver.accept, driver.contentType, driver.auth)
	response := withHeaders.Response()

	return response
}

func (driver *Driver) Post(path string, params hype.Params, data []byte) hype.Response {
	return driver.
		raw.
		Post(path, params, data).
		WithHeaderSet(driver.accept, driver.contentType, driver.auth).
		Response()
}

// PostJSON posts data and returns the status and body whatever the status.
// hype drops the body of a non-2xx response; the validation errors are in it.
func (driver *Driver) PostJSON(path string, data []byte) (int, []byte, error) {
	url := strings.TrimRight(driver.baseURL, "/") + "/" + path
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return 0, nil, err
	}
	for _, header := range []*hype.Header{driver.accept, driver.contentType, driver.auth, driver.userAgent} {
		req.Header.Set(header.Name, header.Value)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return resp.StatusCode, body, err
}
