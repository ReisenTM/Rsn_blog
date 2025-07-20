package ai_service

import (
	"blogX_server/global"
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type ChatResponse struct {
	Id      string `json:"id"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Object  string `json:"object"`
	Usage   struct {
		PromptTokens            int `json:"prompt_tokens"`
		CompletionTokens        int `json:"completion_tokens"`
		TotalTokens             int `json:"total_tokens"`
		CompletionTokensDetails struct {
			AudioTokens     int `json:"audio_tokens"`
			ReasoningTokens int `json:"reasoning_tokens"`
		} `json:"completion_tokens_details"`
		PromptTokensDetails struct {
			AudioTokens  int `json:"audio_tokens"`
			CachedTokens int `json:"cached_tokens"`
		} `json:"prompt_tokens_details"`
	} `json:"usage"`
	SystemFingerprint interface{} `json:"system_fingerprint"`
}

const baseUrl = "https://api.chatanywhere.tech/v1/chat/completions"

func BaseRequest(r Request) (res *http.Response, err error) {
	byteData, _ := json.Marshal(r)
	req, err := http.NewRequest(http.MethodPost, baseUrl, bytes.NewBuffer(byteData))
	if err != nil {
		logrus.Errorf("请求参数失败 %s", err)
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", global.Config.AI.SecretKey))
	req.Header.Add("Content-Type", "application/json")

	res, err = http.DefaultClient.Do(req)
	return
}

//go:embed chat.prompt
var chatPrompt string

// Chat GPT源
func Chat(content string) (msg string, err error) {
	r := Request{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "system",
				Content: chatPrompt,
			},
			{
				Role:    "user",
				Content: content,
			},
		},
	}
	res, err := BaseRequest(r)
	if err != nil {
		return
	}
	body, _ := io.ReadAll(res.Body)
	var response ChatResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		logrus.Errorf("解析失败 %s %s", err, string(body))
		return
	}
	if len(response.Choices) > 0 {
		msg = response.Choices[0].Message.Content
		return
	}
	logrus.Errorf("未获取到数据 %s ", string(body))
	return
}

// DeepSeekUrl DeepSeek源
const DeepSeekUrl = "https://api.deepseek.com/chat/completions"

func DSRequest(r Request) (res *http.Response, err error) {
	byteData, _ := json.Marshal(r)
	req, err := http.NewRequest(http.MethodPost, DeepSeekUrl, bytes.NewBuffer(byteData))
	if err != nil {
		logrus.Errorf("请求参数失败 %s", err)
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", global.Config.AI.SecretKey))
	req.Header.Add("Content-Type", "application/json")

	res, err = http.DefaultClient.Do(req)
	return
}
func DSToChat(content string) (msg string, err error) {
	r := Request{
		Model: "deepseek-chat",
		Messages: []Message{
			{
				Role:    "system",
				Content: chatPrompt,
			},
			{
				Role:    "user",
				Content: content,
			},
		},
	}
	res, err := DSRequest(r)
	if err != nil {
		logrus.Errorln(err)
		return
	}
	body, _ := io.ReadAll(res.Body)
	var response ChatResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		logrus.Errorf("解析失败 %s %s", err, string(body))
		return
	}
	if len(response.Choices) > 0 {
		msg = response.Choices[0].Message.Content
		return
	}
	logrus.Errorf("未获取到数据 %s ", string(body))
	return
}
