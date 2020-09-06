package cmd

import (
	"encoding/json"
	"testing"
)

func TestBuildMessageJson(t *testing.T) {
	type input struct {
		id      string
		message string
	}
	inputs := []input{
		input{
			id:      "1234567890",
			message: "Hello World!",
		},
		// include double quotes
		input{
			id:      "123\"456\"789",
			message: "Hello\" World \"!",
		},
	}
	for _, i := range inputs {
		j := buildMessageJson(i.id, i.message)
		var m map[string]interface{}
		err := json.Unmarshal(j, &m)
		if err != nil {
			t.Errorf("Invalid JSON is built: %v, %v", i.id, i.message)
		}
	}
}
