package freekassa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func sendQuery(url string, payload interface{}) ([]byte, error) {
	var body []byte

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return body, fmt.Errorf("error marshalling payload: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return body, fmt.Errorf("error making HTTP POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return body, fmt.Errorf("error getting good status code: %s", resp.Status)
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return body, fmt.Errorf("error reading response body: %v", err)
	}

	return body, nil
}

func formatMoney(money float64) string {
	return strconv.FormatFloat(money, 'f', -1, 64)
}
