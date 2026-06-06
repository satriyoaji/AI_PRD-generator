package handlers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"prd-generator/internal/models"
	"prd-generator/internal/services"
)

type GenerateHandler struct {
	db        *pgx.Conn
	aiService *services.AIService
}

func NewGenerateHandler(db *pgx.Conn, aiService *services.AIService) *GenerateHandler {
	return &GenerateHandler{db: db, aiService: aiService}
}

// GeneratePRD handles POST /api/generate/prd
func (h *GenerateHandler) GeneratePRD(c *fiber.Ctx) error {
	var req models.GenerateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	transcript, err := h.getTranscript(req.TranscriptID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Transcript not found"})
	}

	content, err := h.aiService.GeneratePRD(transcript.Content)
	if err != nil {
		log.Printf("Failed to generate PRD: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate PRD"})
	}

	generation := models.Generation{
		ID:           uuid.New(),
		TranscriptID: transcript.ID,
		Type:         "prd",
		Content:      json.RawMessage(`{"content": "` + escapeJson(content) + `"}`),
		CreatedAt:    time.Now(),
	}

	h.saveGeneration(generation)

	return c.JSON(generation)
}

// GenerateStories handles POST /api/generate/stories
func (h *GenerateHandler) GenerateStories(c *fiber.Ctx) error {
	var req models.GenerateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	transcript, err := h.getTranscript(req.TranscriptID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Transcript not found"})
	}

	stories, err := h.aiService.GenerateUserStories(transcript.Content)
	if err != nil {
		log.Printf("Failed to generate user stories: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate user stories"})
	}

	content, _ := json.Marshal(map[string]interface{}{"stories": stories})

	generation := models.Generation{
		ID:           uuid.New(),
		TranscriptID: transcript.ID,
		Type:         "stories",
		Content:      content,
		CreatedAt:    time.Now(),
	}

	h.saveGeneration(generation)

	return c.JSON(generation)
}

// GenerateRequirements handles POST /api/generate/requirements
func (h *GenerateHandler) GenerateRequirements(c *fiber.Ctx) error {
	var req models.GenerateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	transcript, err := h.getTranscript(req.TranscriptID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Transcript not found"})
	}

	content, err := h.aiService.GenerateRequirements(transcript.Content)
	if err != nil {
		log.Printf("Failed to generate requirements: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate requirements"})
	}

	generation := models.Generation{
		ID:           uuid.New(),
		TranscriptID: transcript.ID,
		Type:         "requirements",
		Content:      json.RawMessage(`{"content": "` + escapeJson(content) + `"}`),
		CreatedAt:    time.Now(),
	}

	h.saveGeneration(generation)

	return c.JSON(generation)
}

// GenerateWireframe handles POST /api/generate/wireframe
func (h *GenerateHandler) GenerateWireframe(c *fiber.Ctx) error {
	var req models.GenerateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	transcript, err := h.getTranscript(req.TranscriptID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Transcript not found"})
	}

	wireframe, err := h.aiService.GenerateWireframe(transcript.Content)
	if err != nil {
		log.Printf("Failed to generate wireframe: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate wireframe"})
	}

	content, _ := json.Marshal(wireframe)

	generation := models.Generation{
		ID:           uuid.New(),
		TranscriptID: transcript.ID,
		Type:         "wireframe",
		Content:      content,
		CreatedAt:    time.Now(),
	}

	h.saveGeneration(generation)

	return c.JSON(generation)
}

func (h *GenerateHandler) getTranscript(id string) (*models.Transcript, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	if h.db != nil {
		var transcript models.Transcript
		err := h.db.QueryRow(context.Background(),
			"SELECT id, title, content, created_at FROM transcripts WHERE id = $1", parsedID).
			Scan(&transcript.ID, &transcript.Title, &transcript.Content, &transcript.CreatedAt)
		if err != nil {
			return nil, err
		}
		return &transcript, nil
	}

	// Mock response when no DB
	return &models.Transcript{
		ID:        parsedID,
		Title:     "Mock Transcript",
		Content:   "This is a sample meeting transcript about building a new product feature. The team discussed user authentication, dashboard requirements, and reporting features.",
		CreatedAt: time.Now(),
	}, nil
}

func (h *GenerateHandler) saveGeneration(gen models.Generation) {
	if h.db == nil {
		return
	}

	_, err := h.db.Exec(context.Background(),
		"INSERT INTO generations (id, transcript_id, type, content, created_at) VALUES ($1, $2, $3, $4, $5)",
		gen.ID, gen.TranscriptID, gen.Type, gen.Content, gen.CreatedAt)
	if err != nil {
		log.Printf("Failed to save generation: %v", err)
	}
}

func escapeJson(s string) string {
	b, _ := json.Marshal(s)
	return string(b[1 : len(b)-1])
}