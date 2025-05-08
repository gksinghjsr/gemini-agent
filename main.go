package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func main() {
	// Initialize Gemini client with API key from environment variable
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		fmt.Println("GEMINI_API_KEY environment variable is not set")
		return
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}
	defer client.Close()

	scanner := bufio.NewScanner(os.Stdin)
	getUserMessage := func() (string, bool) {
		if !scanner.Scan() {
			return "", false
		}
		return scanner.Text(), true
	}

	tools := []ToolDefinition{
		createCalculatorTool(),
	}
	agent := NewAgent(client, getUserMessage, tools)
	err = agent.Run(ctx)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
}

func NewAgent(client *genai.Client, getUserMessage func() (string, bool), tools []ToolDefinition) *Agent {
	// Convert tools to function declarations
	functionDeclarations := make([]*genai.FunctionDeclaration, len(tools))
	for i, tool := range tools {
		// Define parameters for specific tools
		var parameters *genai.Schema

		if tool.Name == "calculator" {
			parameters = GetCalculatorSchema()
		}

		functionDeclarations[i] = &genai.FunctionDeclaration{
			Name:        tool.Name,
			Description: tool.Description,
			Parameters:  parameters,
		}
	}

	return &Agent{
		client:               client,
		getUserMessage:       getUserMessage,
		tools:                tools,
		functionDeclarations: functionDeclarations,
	}
}

type Agent struct {
	client               *genai.Client
	getUserMessage       func() (string, bool)
	tools                []ToolDefinition
	functionDeclarations []*genai.FunctionDeclaration
}

func (a *Agent) Run(ctx context.Context) error {
	// Create a chat session
	model := a.client.GenerativeModel("gemini-1.5-pro")

	// Set the function declarations as tools
	if len(a.functionDeclarations) > 0 {

		tools := []*genai.Tool{
			{FunctionDeclarations: a.functionDeclarations}}
		model.Tools = tools
	}

	chat := model.StartChat()

	fmt.Println("Chat with Gemini (use 'ctrl-c' to quit)")

	for {
		fmt.Print("You: ")
		userInput, ok := a.getUserMessage()
		if !ok {
			break
		}

		response, err := a.runInference(ctx, chat, userInput)

		if err != nil {
			return err
		}

		fmt.Printf("\u001b[93mGemini\u001b[0m: %v\n", response)
	}

	return nil
}

func (a *Agent) runInference(ctx context.Context, chat *genai.ChatSession, userInput string) (string, error) {
	resp, err := chat.SendMessage(ctx, genai.Text(userInput))
	if err != nil {
		return "", err
	}

	// Check if the response contains function calls
	candidate := resp.Candidates[0]
	if len(candidate.Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from model")
	}

	// Check for function calls in the response
	for _, part := range candidate.Content.Parts {
		if functionCall, ok := part.(genai.FunctionCall); ok {
			// Find the matching tool
			var matchedTool *ToolDefinition
			for i := range a.tools {
				if a.tools[i].Name == functionCall.Name {
					matchedTool = &a.tools[i]
					break
				}
			}

			if matchedTool == nil {
				return "", fmt.Errorf("model called unknown function: %s", functionCall.Name)
			}

			// Convert args to JSON
			argsJSON, err := json.Marshal(functionCall.Args)
			if err != nil {
				return "", fmt.Errorf("error marshaling function args: %v", err)
			}

			// Execute the tool function
			result, err := matchedTool.Function(argsJSON)
			if err != nil {
				return "", fmt.Errorf("error executing function %s: %v", functionCall.Name, err)
			}

			// Send the function result back to the model
			resp, err = chat.SendMessage(ctx, genai.FunctionResponse{
				Name:     functionCall.Name,
				Response: map[string]any{"result": result},
			})
			if err != nil {
				return "", err
			}
		}
	}

	// Process the final text response
	var responseText string
	for _, part := range resp.Candidates[0].Content.Parts {
		if textPart, ok := part.(genai.Text); ok {
			responseText += string(textPart)
		}
	}

	return responseText, nil
}

type ToolDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Function    func(input json.RawMessage) (string, error)
}
