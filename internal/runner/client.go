package runner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nantapop-kj/go-bulkify/config"
)

func ClientRequest(client *http.Client, cfg config.Config, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	req, err := http.NewRequest(cfg.Method, cfg.URL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("new request error: %w", err)
	}

	for k, v := range cfg.Headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	return nil
}
