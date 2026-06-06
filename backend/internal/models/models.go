package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Transcript struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Generation struct {
	ID uuid.UUID       `json:"id"`
	TranscriptID uuid.UUID       `json:"transcript_id"`
	Type         string          `json:"type"` // prd, stories, requirements, wireframe
	Content      json.RawMessage `json:"content"`
	CreatedAt    time.Time       `json:"created_at"`
}

type UserStory struct {
	ID                string   `json:"id"`
	Title             string   `json:"title"`
	Description       string   `json:"description"`
	AcceptanceCriteria []string `json:"acceptance_criteria"`
}

type Wireframe struct {
	MermaidCode string `json:"mermaid_code"`
	Description string `json:"description"`
}

type GenerationResponse struct {
	ID uuid.UUID       `json:"id"`
	TranscriptID uuid.UUID       `json:"transcript_id"`
	Type         string          `json:"type"`
	Content      json.RawMessage `json:"content"`
	CreatedAt    time.Time       `json:"created_at"`
}

type CreateTranscriptRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type GenerateRequest struct {
	TranscriptID string `json:"transcript_id"`
}

type HistoryItem struct {
	ID uuid.UUID `json:"id"`
	TranscriptID uuid.UUID `json:"transcript_id"`
	TranscriptTitle string `json:"transcript_title"`
	Type         string    `json:"type"`
	CreatedAt    time.Time `json:"created_at"`
}
