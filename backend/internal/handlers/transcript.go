package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"path/filepath"
	"strings"
	"time"

	"prd-generator/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

// allowedExtensions lists the only file types accepted by POST /api/transcript.
var allowedExtensions = map[string]bool{
	".txt": true,
	".vtt": true,
}

type Handler struct {
	db *pgx.Conn
}

func NewHandler(db *pgx.Conn) *Handler {
	return &Handler{db: db}
}

// CreateTranscript handles POST /api/transcript
//
// Expects multipart/form-data with:
//   - content  (file, required) — a .txt or .vtt file whose body becomes the transcript text
//   - title    (text, optional) — human-readable title; defaults to the uploaded filename stem
func (h *Handler) CreateTranscript(c *fiber.Ctx) error {
	// --- 1. Retrieve uploaded file from the "content" field ---
	fileHeader, err := c.FormFile("content")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "A file upload is required. Send a .txt or .vtt file in the 'content' multipart field.",
		})
	}

	// --- 2. Validate file extension (.txt or .vtt only) ---
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !allowedExtensions[ext] {
		return c.Status(fiber.StatusUnsupportedMediaType).JSON(fiber.Map{
			"error": "Unsupported file format '" + ext + "'. Only .txt and .vtt files are permitted.",
		})
	}

	// --- 3. Read file bytes and convert to a plain text string ---
	f, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to open uploaded file.",
		})
	}
	defer f.Close()

	raw, err := io.ReadAll(f)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read uploaded file.",
		})
	}

	content := string(raw)
	if strings.TrimSpace(content) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Uploaded file is empty.",
		})
	}

	// --- 4. Derive title: explicit form field > filename stem ---
	title := strings.TrimSpace(c.FormValue("title"))
	if title == "" {
		stem := filepath.Base(fileHeader.Filename)
		title = strings.TrimSuffix(stem, filepath.Ext(stem))
	}

	// --- 5. Persist transcript ---
	transcript := models.Transcript{
		ID:        uuid.New(),
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
	}

	if h.db != nil {
		_, err := h.db.Exec(context.Background(),
			"INSERT INTO transcripts (id, title, content, created_at) VALUES ($1, $2, $3, $4)",
			transcript.ID, transcript.Title, transcript.Content, transcript.CreatedAt)
		if err != nil {
			log.Printf("Failed to save transcript: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to save transcript.",
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(transcript)
}

// GetTranscript handles GET /api/transcript/:id
func (h *Handler) GetTranscript(c *fiber.Ctx) error {
	id := c.Params("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid transcript ID."})
	}

	var transcript models.Transcript
	if h.db != nil {
		err := h.db.QueryRow(context.Background(),
			"SELECT id, title, content, created_at FROM transcripts WHERE id = $1", parsedID).
			Scan(&transcript.ID, &transcript.Title, &transcript.Content, &transcript.CreatedAt)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Transcript not found."})
		}
	} else {
		transcript = models.Transcript{
			ID:        parsedID,
			Title:     "Sample Transcript",
			Content:   "Sample content",
			CreatedAt: time.Now(),
		}
	}

	return c.JSON(transcript)
}

// GetHistory handles GET /api/history
func (h *Handler) GetHistory(c *fiber.Ctx) error {
	var items []models.HistoryItem

	if h.db != nil {
		rows, err := h.db.Query(context.Background(), `
			SELECT g.id, g.transcript_id, COALESCE(t.title, 'Untitled'), g.type, g.created_at
			FROM generations g
			LEFT JOIN transcripts t ON g.transcript_id = t.id
			ORDER BY g.created_at DESC
			LIMIT 50
		`)
		if err != nil {
			log.Printf("Failed to get history: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get history."})
		}
		defer rows.Close()

		for rows.Next() {
			var item models.HistoryItem
			if err := rows.Scan(&item.ID, &item.TranscriptID, &item.TranscriptTitle, &item.Type, &item.CreatedAt); err != nil {
				continue
			}
			items = append(items, item)
		}
	} else {
		items = []models.HistoryItem{}
	}

	return c.JSON(fiber.Map{"items": items})
}

// GetGeneration handles GET /api/history/:id
func (h *Handler) GetGeneration(c *fiber.Ctx) error {
	id := c.Params("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid generation ID."})
	}

	var gen models.Generation
	if h.db != nil {
		err := h.db.QueryRow(context.Background(),
			"SELECT id, transcript_id, type, content, created_at FROM generations WHERE id = $1", parsedID).
			Scan(&gen.ID, &gen.TranscriptID, &gen.Type, &gen.Content, &gen.CreatedAt)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Generation not found."})
		}
	} else {
		gen = models.Generation{
			ID:           parsedID,
			TranscriptID: parsedID,
			Type:         "prd",
			Content:      json.RawMessage(`{"content": "Sample PRD content"}`),
			CreatedAt:    time.Now(),
		}
	}

	return c.JSON(gen)
}
