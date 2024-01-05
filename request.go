package main

import (
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

var (
	systemPrompt = `
		You are a parser of unstructured text data. Your task is to return JSON in the format specified, 
		where each JSON value is filled in using information provided in the data.
		If the data does not contain the information required, return an empty string for that value.

		All data should be returned in CAPITAL LETTERS. There should be no spaces or punctuation in any field EXCEPT for names.

		For example: if the field is 08-23-87, return 082387. If the field is 08/23/87, return 082387.

		If you are unsure, or unable to parse a field, simply return an empty string.
	`
	jsonTagPrompt = `
		[JSON Structure]
	`

	dataTagPrompt = `
		[Data to fill JSON structure with]
	`

	filledJSONTagPrompt = `
		[Filled JSON structure]
	`
	openAIKey = os.Getenv("OPENAI_API_KEY")
	model     = openai.GPT4
)

func requestFill(jsonStructure string, inputData string) (string, error) {
	var messages = []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		},
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: fmt.Sprintf(`%s`, jsonTagPrompt),
		},
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: fmt.Sprintf(`%s`, jsonStructure),
		},
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: fmt.Sprintf(`%s`, dataTagPrompt),
		},
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: fmt.Sprintf(`%s`, inputData),
		},
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: fmt.Sprintf(`%s`, filledJSONTagPrompt),
		},
	}

	if openAIKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY is not set")
	}

	client := openai.NewClient(openAIKey)

	response, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    model,
			Messages: messages,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to create chat completion: %w", err)
	}

	completion := response.Choices[0].Message.Content

	return completion, nil
}
