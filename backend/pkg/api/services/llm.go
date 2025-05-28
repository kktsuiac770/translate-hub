package services

import (
	"context"
	"log"
	"os"

	genai "google.golang.org/genai"
)

// MockLLMTranslate simulates LLM translation for MVP
func MockLLMTranslate(text string) string {
	return "[LLM] " + text // prepend for demo
}

// GeminiTranslate uses Google GenAI to translate text from sourceLang to targetLang
func GeminiTranslate(text, sourceLang, targetLang string) string {
	apiKey := os.Getenv("GOOGLE_GENAI_API_KEY")
	if apiKey == "" {
		return "[LLM error] missing API key"
	}
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatal(err)
	}
	prompt := "Translate the following sentence from " + sourceLang + " to " + targetLang + ". Output only the translated sentence, no explanation, no extra text.\n" + text
	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		genai.Text(prompt),
		nil,
	)

	if err != nil {
		return "[LLM error] " + text
	}

	return result.Text()
}
