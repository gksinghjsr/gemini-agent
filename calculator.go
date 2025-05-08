package main

import (
	"encoding/json"
	"fmt"

	"github.com/google/generative-ai-go/genai"
)

// Define a struct for calculator input parameters
type CalculatorInput struct {
	Operation string  `json:"operation"`
	A         float64 `json:"a"`
	B         float64 `json:"b"`
}

// Calculator tool function
func calculatorTool(input json.RawMessage) (string, error) {
	var params CalculatorInput
	fmt.Println("\n[DEBUG] Calculator tool called with input:", string(input))
	if err := json.Unmarshal(input, &params); err != nil {
		return "", fmt.Errorf("failed to parse calculator parameters: %v", err)
	}

	var result float64
	switch params.Operation {
	case "add":
		result = params.A + params.B
	case "subtract":
		result = params.A - params.B
	case "multiply":
		result = params.A * params.B
	case "divide":
		if params.B == 0 {
			return "", fmt.Errorf("division by zero")
		}
		result = params.A / params.B
	default:
		return "", fmt.Errorf("unknown operation: %s", params.Operation)
	}
	return fmt.Sprintf("%.2f", result), nil
}

func createCalculatorTool() ToolDefinition {
	return ToolDefinition{
		Name:        "calculator",
		Description: "Performs basic arithmetic operations: add, subtract, multiply, divide",
		Function:    calculatorTool,
	}
}

// GetCalculatorSchema returns the schema definition for the calculator tool
func GetCalculatorSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"operation": {
				Type:        genai.TypeString,
				Enum:        []string{"add", "subtract", "multiply", "divide"},
				Description: "The arithmetic operation to perform",
			},
			"a": {
				Type:        genai.TypeNumber,
				Description: "The first operand",
			},
			"b": {
				Type:        genai.TypeNumber,
				Description: "The second operand",
			},
		},
		Required: []string{"operation", "a", "b"},
	}
}
