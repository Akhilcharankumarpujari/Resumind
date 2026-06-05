package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)



const systemPrompt = `You are an expert ATS (Applicant Tracking System) specialist and career coach with 15+ years of experience reviewing resumes for top-tier companies.

Your task is to analyze the provided resume text against the job title and job description, and return a detailed, honest evaluation.

IMPORTANT RULES:
- Be thorough and critical. Do NOT give artificially high scores.
- Score honestly based on actual resume quality vs. the job requirements.
- A poor resume for a role SHOULD receive a low score (e.g. 30-50).
- Identify missing keywords, weak action verbs, poor formatting signals (from text structure), and skill gaps.
- Each section must have 3-4 tips minimum.
- Tips with type "good" praise something done well. Tips with type "improve" point out specific problems.
- The "tip" field is a short headline. The "explanation" field is a detailed explanation.

Return ONLY a valid JSON object matching this exact structure, with no extra text, no markdown backticks:

{
  "overallScore": <number 0-100>,
  "ATS": {
    "score": <number 0-100>,
    "tips": [
      { "type": "good"|"improve", "tip": "<short headline>", "explanation": "<detailed explanation>" }
    ]
  },
  "toneAndStyle": {
    "score": <number 0-100>,
    "tips": [
      { "type": "good"|"improve", "tip": "<short headline>", "explanation": "<detailed explanation>" }
    ]
  },
  "content": {
    "score": <number 0-100>,
    "tips": [
      { "type": "good"|"improve", "tip": "<short headline>", "explanation": "<detailed explanation>" }
    ]
  },
  "structure": {
    "score": <number 0-100>,
    "tips": [
      { "type": "good"|"improve", "tip": "<short headline>", "explanation": "<detailed explanation>" }
    ]
  },
  "skills": {
    "score": <number 0-100>,
    "tips": [
      { "type": "good"|"improve", "tip": "<short headline>", "explanation": "<detailed explanation>" }
    ]
  }
}`


func AnalyzeResume(resumeText, jobTitle, jobDesc string) (*Feedback, error) {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GROQ_API_KEY environment variable is not set")
	}


	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://api.groq.com/openai/v1"
	client := openai.NewClientWithConfig(config)

	userMessage := fmt.Sprintf(
		"Job Title: %s\n\nJob Description:\n%s\n\nResume Text:\n%s",
		jobTitle, jobDesc, resumeText,
	)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "llama-3.3-70b-versatile",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userMessage,
				},
			},

			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONObject,
			},
			Temperature: 0.3,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("Groq API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("Groq returned no response choices")
	}

	rawJSON := resp.Choices[0].Message.Content

	var feedback Feedback
	if err := json.Unmarshal([]byte(rawJSON), &feedback); err != nil {
		return nil, fmt.Errorf("failed to parse AI JSON response: %w\nRaw response: %s", err, rawJSON)
	}

	return &feedback, nil
}
