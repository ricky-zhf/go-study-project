package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("You:")
		scanner.Scan()
		input := scanner.Text()
		if input == "exit" {
			break
		}
		output, err := sendGPTRequest(input)
		if err != nil {
			fmt.Printf("Failed. %v\n", err)
			continue
		}

		fmt.Println("GPT:", output)
	}
}

var client = &http.Client{}

type Req struct {
	Model string     `json:"model"`
	Msg   []Messages `json:"messages"`
}

type Messages struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func sendGPTRequest(inputText string) (string, error) {
	//url := "https://api.openai.com/v1/chat/completions"
	var (
		response ApifoxModel
		url      = "https://api.chatanywhere.com.cn/v1/chat/completions"
		api_key  = "sk-LfBYwIIyuExut9xUkBhksU2vghyAhoe7npPGrPIILeGl8zim"
	)

	msg := Messages{
		Role:    "user",
		Content: inputText,
	}
	payload, err := json.Marshal(Req{
		Model: "gpt-3.5-turbo",
		Msg:   []Messages{msg},
	})
	if err != nil {
		return "", fmt.Errorf("marshal failed. %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return "", fmt.Errorf("new request failed. %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+api_key)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("req failed. %v", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read failed. %v", err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("unmarshal failed. %v", err)
	}

	return response.Choices[0].Message.Content, nil
}

type ApifoxModel struct {
	Choices []Choice `json:"choices"`
	Created int64    `json:"created"`
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	FinishReason *string  `json:"finish_reason,omitempty"`
	Index        *int64   `json:"index,omitempty"`
	Message      *Message `json:"message,omitempty"`
}

type Message struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type Usage struct {
	CompletionTokens int64 `json:"completion_tokens"`
	PromptTokens     int64 `json:"prompt_tokens"`
	TotalTokens      int64 `json:"total_tokens"`
}
