package http

import (
	"bytes"
	"encoding/json"
	"golang.org/x/net/html/charset"
	"io"
	"net/http"
	"time"
)

func PostJson(url string, jsonData interface{}, headers map[string]string) ([]byte, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	jsonBody, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyReader, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(bodyReader)
	if err != nil {
		return nil, err
	}

	return data, nil
}
