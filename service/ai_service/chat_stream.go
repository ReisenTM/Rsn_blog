// 流式输出
package ai_service

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"github.com/sirupsen/logrus"
)

type Choice struct {
	Index int `json:"index"`
	Delta struct {
		Content string `json:"content"`
	} `json:"delta"`
	Logprobs     interface{} `json:"logprobs"`
	FinishReason interface{} `json:"finish_reason"`
}

type StreamData struct {
	Id                string      `json:"id"`
	Choices           []Choice    `json:"choices"`
	Created           int         `json:"created"`
	Model             string      `json:"model"`
	Object            string      `json:"object"`
	SystemFingerprint interface{} `json:"system_fingerprint"`
}

//go:embed chat_stream.prompt
var chatStreamPrompt string

func ChatStream(content string, params string) (msgChan chan string, err error) {
	msgChan = make(chan string)
	r := Request{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "system",
				Content: chatStreamPrompt + params,
			},
			{
				Role:    "user",
				Content: content,
			},
		},
		Stream: true,
	}
	res, err := BaseRequest(r)
	if err != nil {
		return
	}
	//逐行读取内容，适合持续推送的场景
	scanner := bufio.NewScanner(res.Body)
	// 按行分割
	scanner.Split(bufio.ScanLines)
	go func() {
		for scanner.Scan() {
			text := scanner.Text()
			if text == "" {
				continue
			}
			//把前缀 "data: " 去掉
			data := text[6:]

			var item StreamData
			err = json.Unmarshal([]byte(data), &item)
			if err != nil {
				logrus.Errorf("解析失败 %s %s", err, data)
				continue

			}
			if len(item.Choices) <= 0 {
				close(msgChan)
				return
			}
			msgChan <- item.Choices[0].Delta.Content
		}
	}()

	return
}

// ChatStream DS源
func DSChatStream(content string, params string) (msgChan chan string, err error) {
	msgChan = make(chan string)
	r := Request{
		Model: "deepseek-chat",
		Messages: []Message{
			{
				Role:    "system",
				Content: chatStreamPrompt + params,
			},
			{
				Role:    "user",
				Content: content,
			},
		},
		Stream: true,
	}
	res, err := DSRequest(r)
	if err != nil {
		return
	}
	//逐行读取内容，适合持续推送的场景
	scanner := bufio.NewScanner(res.Body)
	// 按行分割
	scanner.Split(bufio.ScanLines)
	go func() {
		for scanner.Scan() {
			text := scanner.Text()
			if text == "" {
				continue
			}
			//把前缀 "data: " 去掉
			data := text[6:]
			if data == "[DONE]" {
				close(msgChan)
				return
			}
			var item StreamData
			err = json.Unmarshal([]byte(data), &item)
			if err != nil {
				logrus.Errorf("解析失败 %s %s", err, data)
				continue

			}
			if len(item.Choices) <= 0 {
				close(msgChan)
				return
			}
			msgChan <- item.Choices[0].Delta.Content
		}
	}()

	return
}
