package main

import (
	"encoding/json"

	"github.com/invopop/jsonschema"
)

func decodeData(inputData string, dataStructure any) (any, error) {
	jsonFormat := jsonschema.Reflect(dataStructure)

	schemaBytes, err := jsonFormat.MarshalJSON()
	if err != nil {
		return nil, err
	}

	schemaString := string(schemaBytes)

	completion, err := requestFill(schemaString, inputData)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(completion), &dataStructure)
	if err != nil {
		return err, nil
	}

	return dataStructure, nil
}
