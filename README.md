# Gemini Agent

A Go application that creates an interactive chat interface with Google's Gemini AI model, enhanced with a calculator tool.

## Getting started

1. Clone the repo
   ```bash
   git clone https://github.com/gksinghjsr/gemini-agent.git
   ```

2. Navigate to the project directory
   ```bash
   cd gemini-agent
   ```

3. Set your Gemini API key as an environment variable
   ```bash
   export GEMINI_API_KEY=your_api_key_here
   ```

4. Build and run the application
   ```bash
   go run main.go calculator.go
   ```

## Usage

Once running, you'll see a prompt: "Chat with Gemini (use 'ctrl-c' to quit)"

You can:
- Type messages to interact with Gemini
- Ask it to perform calculations using the calculator tool
  - Example: "Calculate 25 multiplied by 4"
  - The model will use the calculator tool with operations: add, subtract, multiply, or divide

To exit the application, press Ctrl+C.

## Features

- Interactive chat interface with Gemini 1.5 Pro
- Calculator tool for basic arithmetic operations (add, subtract, multiply, divide)
- Function calling capabilities to demonstrate AI agent functionality