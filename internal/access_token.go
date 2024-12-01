package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

// Функция для создания HTTP-запроса
func createRequest(method, url string, payload []byte, headers map[string]string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return req, nil
}

// Функция для отправки HTTP-запроса
func sendRequest(req *http.Request) (*http.Response, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	return resp, nil
}

// Функция для обработки ответа
func handleResponse(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, nil
}

// Функция для извлечения токена доступа из ответа
func extractAccessToken(body []byte) (string, error) {
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	if accessToken, ok := result["access_token"].(string); ok {
		return accessToken, nil
	}

	return "", fmt.Errorf("access token not found in response")
}

// Функция для извлечения ответа от нейросети из ответа
func extractNeuralNetResponse(body []byte) (string, error) {
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	choices, err := getChoices(result)
	if err != nil {
		return "", err
	}

	for _, choice := range choices {
		content, err := getContent(choice)
		if err != nil {
			return "", err
		}
		if content != "" {
			return content, nil
		}
	}

	return "", fmt.Errorf("neural net response not found in response")
}

// Функция для извлечения choices из ответа
func getChoices(result map[string]interface{}) ([]interface{}, error) {
	choices, ok := result["choices"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("choices not found in response")
	}
	return choices, nil
}

// Функция для извлечения content из choice
func getContent(choice interface{}) (string, error) {
	choiceMap, ok := choice.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid choice format")
	}

	message, ok := choiceMap["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("message not found in choice")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("content not found in message")
	}

	return content, nil
}

// Функция для получения токена доступа
func getAccessToken(authorizationKey string) (string, error) {
	rqUID := uuid.New().String()
	url := "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"
	payload := "scope=GIGACHAT_API_PERS"

	headers := map[string]string{
		"Content-Type":  "application/x-www-form-urlencoded",
		"Accept":        "application/json",
		"RqUID":         rqUID,
		"Authorization": "Basic " + authorizationKey,
	}

	req, err := createRequest("POST", url, []byte(payload), headers)
	if err != nil {
		return "", err
	}

	resp, err := sendRequest(req)
	if err != nil {
		return "", err
	}

	body, err := handleResponse(resp)
	if err != nil {
		return "", err
	}

	return extractAccessToken(body)
}

// Функция для отправки запроса к нейросети
func sendNeuralNetRequest(accessToken string) (string, error) {
	url := "https://gigachat.devices.sberbank.ru/api/v1/chat/completions"
	payload := map[string]interface{}{
		"model": "GigaChat",
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": "Привет! Как дела?",
			},
		},
		"stream":             false,
		"repetition_penalty": 1,
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Accept":        "application/json",
		"Authorization": "Bearer " + accessToken,
	}

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error marshalling request body: %w", err)
	}

	req, err := createRequest("POST", url, reqBody, headers)
	if err != nil {
		return "", err
	}

	resp, err := sendRequest(req)
	if err != nil {
		return "", err
	}

	body, err := handleResponse(resp)
	if err != nil {
		return "", err
	}

	return extractNeuralNetResponse(body)
}

func main() {
	authorizationKey := "OGYwMTljNzYtYzEyMy00MjE4LWJmY2UtZTY2ZWE1ZGRlM2E4Ojg1ZjcyYWE4LWNmMTgtNGIyMC05NzQ5LTE5Njk4ZTk1MmExMA=="

	accessToken, err := getAccessToken(authorizationKey)
	if err != nil {
		fmt.Println("Error getting access token:", err)
		return
	}
	fmt.Println("Access Token:", accessToken)

	neuralNetResponse, err := sendNeuralNetRequest(accessToken)
	if err != nil {
		fmt.Println("Error sending neural net request:", err)
		return
	}
	fmt.Println("Neural Net Response Content:", neuralNetResponse)
}
