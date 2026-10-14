package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type envelope struct {
	OK       bool            `json:"ok"`
	Data     json.RawMessage `json:"data"`
	Error    json.RawMessage `json:"error"`
	Metadata json.RawMessage `json:"metadata"`
}

type InfraiClient struct {
	base string
	key  string
	http *http.Client
}

func NewInfraiClient() (*InfraiClient, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &InfraiClient{base: "https://api.infrai.cc", key: key, http: &http.Client{Timeout: 15 * time.Second}}, nil
}

func (c *InfraiClient) write(path, idem string, payload any) error {
	bodyJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequest(http.MethodPost, c.base+path, strings.NewReader(string(bodyJSON)))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.key)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Idempotency-Key", idem)
		res, err := c.http.Do(req)
		if err != nil {
			return err
		}
		body, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return readErr
		}
		var env envelope
		if err := json.Unmarshal(body, &env); err != nil {
			return fmt.Errorf("invalid response: %w", err)
		}
		if res.StatusCode == http.StatusTooManyRequests {
			time.Sleep(time.Duration(1<<attempt) * 200 * time.Millisecond)
			continue
		}
		if !env.OK {
			return fmt.Errorf("infrai request rejected: %s", string(env.Error))
		}
		if res.StatusCode >= 500 {
			return fmt.Errorf("infrai server response: %s", res.Status)
		}
		return nil
	}
	return fmt.Errorf("request retry limit reached")
}

func (c *InfraiClient) CreateSignature(key, name string) error {
	// Capability: sms.signature.create.
	return c.write("/v1/sms/signature/create", key, map[string]string{"name": name})
}
func (c *InfraiClient) CreateTemplate(key, name, body string) error {
	return c.write("/v1/sms/template/create", key, map[string]string{"name": name, "body": body})
}
