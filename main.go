package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/invopop/jsonschema"
	"github.com/sashabaranov/go-openai"
)

// Struct - Defined by user
type Address struct {
	Street string
	City   string
	State  string
}

type ReturnData struct {
	Name    string
	Age     int
	Address Address
}

func main() {
	var userType ReturnData
	decodeData(`- Name: John Smith - Age: 30 - Address: 123 Main St, New York, NY`, &userType)

}

func decodeData(inputData string, dataStructure any) error {
	// Create JSON schema from struct
	jsonFormat := jsonschema.Reflect(dataStructure)
	// Get JSON Schema in string format
	jsonFormatInString, err := jsonFormat.MarshalJSON()
	if err != nil {
		return err
	}

	client := openai.NewClient(os.Getenv("OPEN_API_KEY"))
	response, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: `
					You are a parser unstructured text data. I need you to return raw JSON in your answers.
					This JSON will need to be in the format described in the message- and need to
					`,
				},
				{
					Role: openai.ChatMessageRoleUser,
					Content: fmt.Sprintf(`
						%s
						I'd like to format the following data in the above JSON string. It has to be returned like this as it's going to be parsed by code.
						The data is:
                        %s
				`, jsonFormatInString, inputData),
				},
			},
		},
	)
	if err != nil {
		panic(err)
	}
	// Get the first message choice
	completion := response.Choices[0].Message.Content

	// Convert AI generated JSON to out struct
	err = json.Unmarshal([]byte(completion), &dataStructure)

	if err != nil {
		// PANIK
		panic(err)
	}

	return nil
}
