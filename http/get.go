package http

import (
	"github.com/corpix/uarand"
	"golang.org/x/net/html/charset"
	"io"
	"net/http"
	"time"
)

func Get(url string, headers map[string]interface{}, timeout int) ([]byte, error, string) {
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	req, _ := http.NewRequest("GET", url, nil)
	for k, v := range headers {
		req.Header.Set(k, v.(string))
	}
	req.Header.Set("User-Agent", uarand.GetRandom())

	resp, err := client.Do(req)
	if err != nil {
		return nil, err, ""
	}
	defer resp.Body.Close()
	bodyReader, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, err, ""
	}

	data, err := io.ReadAll(bodyReader)
	if err != nil {
		return nil, err, ""
	}

	return data, nil, resp.Header.Get("Content-Type")
}
