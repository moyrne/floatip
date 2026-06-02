package detector

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func DetectPublicIP(url string) (string, error) {
	if url == "" {
		return "", fmt.Errorf("IP echo URL is empty")
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to query IP echo service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("IP echo service returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read IP echo response: %w", err)
	}

	ip := strings.TrimSpace(string(body))
	if ip == "" {
		return "", fmt.Errorf("IP echo service returned empty body")
	}

	return ip, nil
}
