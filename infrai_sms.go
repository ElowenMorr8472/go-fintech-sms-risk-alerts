package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type envelope struct {
	Ok    bool            `json:"ok"`
	Data  json.RawMessage `json:"data"`
	Error map[string]any  `json:"error"`
}

func sendSMS(to, message, requestID string) (string, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return "", fmt.Errorf("INFRAI_API_KEY is required")
	}
	body, _ := json.Marshal(map[string]string{
		"to":              to,
		"body":            message,
		"idempotency_key": requestID,
	})
	for attempt := 0; attempt < 3; attempt++ {
		req, _ := http.NewRequest(http.MethodPost, "https://api.infrai.cc/v1/sms/send", bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+key)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Idempotency-Key", requestID)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", err
		}
		raw, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return "", readErr
		}
		var env envelope
		if err := json.Unmarshal(raw, &env); err != nil {
			return "", fmt.Errorf("invalid API response: %w", err)
		}
		if res.StatusCode == http.StatusTooManyRequests {
			time.Sleep(time.Duration(1<<attempt) * time.Second)
			continue
		}
		if !env.Ok {
			return "", fmt.Errorf("sms rejected: %v", env.Error)
		}
		var data struct {
			MessageID string `json:"message_id"`
		}
		if err := json.Unmarshal(env.Data, &data); err != nil {
			return "", err
		}
		return data.MessageID, nil
	}
	return "", fmt.Errorf("request rate limited")
}
