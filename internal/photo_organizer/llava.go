package photoorganizer

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	model       = "llava"
	prompt      = "Generate a JSON response with the fields 'description' (a brief image description under 20 words) and 'tags' (a list of relevant keywords summarizing the image's subjects, objects, or themes)."
	ollamaApi   = "http://localhost:11434/api/generate"
	temperature = 0
)

type OllamaRequest struct {
	Model       string   `json:"model"`
	Prompt      string   `json:"prompt"`
	Images      []string `json:"images"`
	Stream      bool     `json:"stream"`
	Temperature float32  `json:"temperature"`
}

type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

type LlavaDescriber struct{}

func (d LlavaDescriber) Describe(imagePath string) (string, []string, error) {
	return LlavaDescribeImage(imagePath)
}

func LlavaDescribeImage(imagePath string) (string, []string, error) {
	imgBase64, err := encodeImageToBase64(imagePath)
	if err != nil {
		return "", nil, fmt.Errorf("error encoding image: %v", err)
	}

	reqData := OllamaRequest{
		Model:       model,
		Prompt:      prompt,
		Images:      []string{imgBase64},
		Stream:      false,
		Temperature: temperature,
	}

	body, _ := json.Marshal(reqData)
	resp, err := http.Post(ollamaApi, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", nil, fmt.Errorf("error on Post: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, fmt.Errorf("error reading the response body: %v", err)
	}

	var response OllamaResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return "", nil, fmt.Errorf("error decoding the response: %v", err)
	}

	return extractMetadata(response.Response)
}

type ImageMetadata struct {
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

func extractMetadata(response string) (string, []string, error) {
	jsonData := strings.TrimSpace(strings.ReplaceAll(response, "```json", ""))
	jsonData = strings.TrimSpace(strings.ReplaceAll(jsonData, "```", ""))

	var metadata ImageMetadata
	err := json.Unmarshal([]byte(jsonData), &metadata)
	if err != nil {
		fmt.Println(jsonData)
		fmt.Println("Error al decodificar JSON:", err)
		return "", nil, err
	}

	return metadata.Description, metadata.Tags, nil
}

func encodeImageToBase64(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}
