package main

import (
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

var (
	systemPrompt = `
		You are a parser unstructured text data. I need you to return raw JSON in your answers.
		This JSON will need to be in the format described in the message
	`
	userPrompt = `
		I'd like to format the following data in the above JSON string. It has to be returned like this as it's going to be parsed by code.
		Ths data is:
	`
	openAIKey = os.Getenv("OPEN_AI_KEY")
	model = openai.GPT4
)

func requestFill(format string, inputData string) (string, error) {
	var messages = []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		},
		{
			Role: openai.ChatMessageRoleUser,
			Content: fmt.Sprintf(`%s\n%s\n%s`, format, userPrompt, inputData),
		},
	}

	client := openai.NewClient(openAIKey)

	response, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
			Messages: messages,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to create chat completion: %w", err)
	}

	completion := response.Choices[0].Message.Content

	return completion, nil
}
