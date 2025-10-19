package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rojack96/endoflife-bot/config"
	"go.uber.org/zap"
)

// HTTPRequest is a generic function to make HTTP requests and return a json response
func HttpRequest[T any](method, url string, body any, result *T) error {
	var reqBody io.Reader

	log := config.GetLogger()

	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			log.Error("error marshalling request body", zap.Error(err))
			return fmt.Errorf("error an marshalling request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		log.Error("error creating HTTP request", zap.Error(err))
		return fmt.Errorf("error creating HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("error making HTTP request", zap.Error(err))
		return fmt.Errorf("error making HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Error("received non-2xx response", zap.Int("status_code", resp.StatusCode), zap.String("body", string(bodyBytes)))
		return fmt.Errorf("error HTTP %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		log.Error("error decoding JSON response", zap.Error(err))
		return fmt.Errorf("errore nel decode JSON: %w", err)
	}

	return nil
}
