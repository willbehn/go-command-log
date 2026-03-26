package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
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

type Response struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
}

type Result struct {
	Summary  string   `json:"summary"`
	Commands []string `json:"commands"`
	Risk     string   `json:"risk"`
	Notes    []string `json:"notes"`
}

var (
	model    string
	endpoint string
)

var askCmd = &cobra.Command{
	Use:   "ask",
	Short: "Ask for shell command guidance quickly from your terminal",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		task := strings.Join(args, " ")
		prompt := buildPrompt(task)

		responseText, err := callModel(endpoint, model, prompt)
		if err != nil {
			return err
		}

		var result Result
		if err := json.Unmarshal([]byte(responseText), &result); err != nil {
			return fmt.Errorf("todo %w", err)
		}

		printResult(result)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(askCmd)

	askCmd.Flags().StringVar(&model, "model", "mistral", "Model name for the local API")
	askCmd.Flags().StringVar(&endpoint, "endpoint", "http://localhost:11434/api/chat", "Chat API endpoint")
}

func callModel(endpoint, model, prompt string) (string, error) {
	reqBody := Request{
		Model:    model,
		Messages: []Message{{Role: "user", Content: prompt}},
		Stream:   false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("model API returned status %s: %s", resp.Status, string(body))
	}

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("decode response envelope: %w", err)
	}

	return result.Message.Content, nil
}

func printResult(result Result) {
	fmt.Printf("Summary: %s\n", result.Summary)
	fmt.Printf("Risk: %s\n", result.Risk)

	if len(result.Commands) > 0 {
		fmt.Println("Commands:")
		for i, c := range result.Commands {
			fmt.Printf("  %d. %s\n", i+1, c)
		}
	}

	if len(result.Notes) > 0 {
		fmt.Println("Notes:")
		for i, n := range result.Notes {
			fmt.Printf("  %d. %s\n", i+1, n)
		}
	}
}

func buildPrompt(userInput string) string {
	return `
		You are a Unix command-line assistant.

		Return ONLY valid JSON. No explanations.

		Format:
		{
		"summary": string,
		"commands": string[],
		"risk": "safe" | "caution" | "dangerous",
		"notes": string[]
		}

		Rules:

		1. Commands must be valid bash commands.
		2. Prefer SAFE and ROBUST commands over short ones.
		3. NEVER use placeholders like /path/to/... if the user gave a concrete name.
		4. If the task is destructive (delete, overwrite, move many files):
		- Include a SAFE preview command first (e.g. using find without -delete)
		- Then include the actual execution command
		5. If the request is ambiguous, DO NOT guess:
		- Return an empty "commands" array
		- Add a note asking for clarification
		6. Risk levels:
		- "safe": read-only operations (ls, grep, find without delete)
		- "caution": file modifications (cp, mv, chmod, non-recursive operations)
		- "dangerous": destructive operations (rm, -delete, overwriting many files)
		7. Prefer precise tools over fragile ones:
		- Prefer "find" over "ls | awk"
		- Handle edge cases like spaces and case sensitivity when reasonable
		8. Preserve user intent exactly (filenames, directories, wording).

		Task:
		` + userInput
}
