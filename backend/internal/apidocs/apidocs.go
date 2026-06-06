package apidocs

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type Example struct {
	Language string
	Value    string
}

type Endpoint struct {
	Method              string
	Path                string
	Summary             string
	Description         string
	Tags                []string
	RequestContentType  string
	RequestDescription  string
	RequestExample      string
	ResponseDescription string
	ResponseExample     string
	ErrorResponses      map[string]string
}

var endpoints = []Endpoint{
	{
		Method:              "GET",
		Path:                "/health",
		Summary:             "Health check",
		Description:         "Returns the API health status.",
		Tags:                []string{"Health"},
		ResponseDescription: "Service health status.",
		ResponseExample: `{
  "status": "healthy",
  "service": "prd-generator-api"
}`,
	},
	{
		Method:             "POST",
		Path:               "/api/transcript",
		Summary:            "Create transcript from uploaded file",
		Description:        "Creates a transcript by uploading a plain text `.txt` file or WebVTT `.vtt` file. The uploaded file must be sent as multipart form-data using the `content` file field. The file body is read and stored as a text string.",
		Tags:               []string{"Transcripts"},
		RequestContentType: "multipart/form-data",
		RequestDescription: "Multipart form with optional `title` text field and required `content` file field. Only `.txt` and `.vtt` file extensions are permitted.",
		RequestExample: `curl -X POST http://localhost:8080/api/transcript \
  -F "title=Q1 Planning Meeting" \
  -F "content=@meeting.vtt"`,
		ResponseDescription: "Created transcript. `content` contains the uploaded file contents converted to a string.",
		ResponseExample: `{
  "id": "2f4fb8e8-8b51-47d0-b699-3fc2d64c6f2d",
  "title": "Q1 Planning Meeting",
  "content": "WEBVTT\n\n00:00:00.000 --> 00:00:05.000\nWe need a dashboard for product managers.",
  "created_at": "2026-06-06T10:30:00Z"
}`,
		ErrorResponses: map[string]string{
			"400": "Missing file, empty file, or malformed multipart request.",
			"415": "Uploaded file extension is not `.txt` or `.vtt`.",
			"500": "Failed to read or save transcript.",
		},
	},
	{
		Method:              "GET",
		Path:                "/api/transcript/{id}",
		Summary:             "Get transcript by ID",
		Description:         "Returns a previously created transcript.",
		Tags:                []string{"Transcripts"},
		ResponseDescription: "Transcript details.",
		ResponseExample: `{
  "id": "2f4fb8e8-8b51-47d0-b699-3fc2d64c6f2d",
  "title": "Q1 Planning Meeting",
  "content": "We need a dashboard for product managers.",
  "created_at": "2026-06-06T10:30:00Z"
}`,
		ErrorResponses: map[string]string{
			"400": "Invalid transcript ID.",
			"404": "Transcript not found.",
		},
	},
	{
		Method:             "POST",
		Path:               "/api/generate/prd",
		Summary:            "Generate PRD",
		Description:        "Generates a Product Requirements Document from a transcript.",
		Tags:               []string{"Generation"},
		RequestContentType: "application/json",
		RequestDescription: "JSON body containing the transcript ID.",
		RequestExample: `{
  "transcript_id": "2f4fb8e8-8b51-47d0-b699-3fc2d64c6f2d"
}`,
		ResponseDescription: "Generated PRD content.",
		ResponseExample: `{
  "id": "b6a0ce5a-a43f-414f-83b4-8a908739bcfb",
  "transcript_id": "2f4fb8e8-8b51-47d0-b699-3fc2d64c6f2d",
  "type": "prd",
  "content": {
    "content": "# Product Requirements Document\n\n## Overview\nBuild a dashboard for product managers."
  },
  "created_at": "2026-06-06T10:31:00Z"
}`,
		ErrorResponses: generationErrors("PRD"),
	},
	{
		Method:             "POST",
		Path:               "/api/generate/stories",
		Summary:            "Generate user stories",
		Description:        "Generates user stories from a transcript.",
		Tags:               []string{"Generation"},
		RequestContentType: "application/json",
		RequestDescription: "JSON body containing the transcript ID.",
		RequestExample: `{
  "transcript_id": "2f4fb8e8-8b51-47d0-b699-3fc2d64c6f2d"
}`,
		ResponseDescription: "Generated user stories.",
		ResponseExample: `{
  "id": "6a2f4c1a-bd1d-4b8a-82ed-09c0ced1dc3d",
  "transcript_id": "2f4fb8e8-8b51-47d0-b699-3fc2d64c6f2d",
  "type": "stories",
  "content": {
    "stories": [
      {
        "id": "US-001",
        "title": "View dashboard metrics",
        "description": "As a product manager, I want to view dashboard metrics so I can track product performance.",
        "acceptance_criteria": ["Dashboard displays key metrics", "Data refreshes successfully"]
      }
    ]
  },
  "created_at": "2026-06-06T10:32:00Z"
}`,
		ErrorResponses: generationErrors("user stories"),
	},
	{
		Method:             "POST",
		Path:               "/api/generate/requirements",
		Summary:            "Generate functional requirements",
		Description:        "Generates functional requirements from a transcript.",
		Tags:               []string{"Generation"},
		RequestContentType: "application/json",
		RequestDescription: "JSON body containing the transcript ID.",
		RequestExample: `{
  "transcript_id": "2f4fb8e8-8b51-47d0-b699-3fc2d64c6f2d"
}`,
		ResponseDescription: "Generated functional requirements.",
		ResponseExample: `{
  "id": "a260d2d4-d2aa-4e55-8765-51d1530748af",
  "transcript_id": "2f4fb8e8-8b51-47d0-b699-3fc2d64c6f2d",
  "type": "requirements",
  "content": {
    "content": "# Functional Requirements\n\n1. The system shall display dashboard metrics."
  },
  "created_at": "2026-06-06T10:33:00Z"
}`,
		ErrorResponses: generationErrors("requirements"),
	},
	{
		Method:             "POST",
		Path:               "/api/generate/wireframe",
		Summary:            "Generate wireframe",
		Description:        "Generates UI/UX wireframe documentation from a transcript.",
		Tags:               []string{"Generation"},
		RequestContentType: "application/json",
		RequestDescription: "JSON body containing the transcript ID.",
		RequestExample: `{
  "transcript_id": "2f4fb8e8-8b51-47d0-b699-3fc2d64c6f2d"
}`,
		ResponseDescription: "Generated wireframe content.",
		ResponseExample: `{
  "id": "6957633d-7716-46b7-8bd3-3d0d87e8a318",
  "transcript_id": "2f4fb8e8-8b51-47d0-b699-3fc2d64c6f2d",
  "type": "wireframe",
  "content": {
    "mermaid_code": "graph TD\n  A[Dashboard] --> B[Metrics Cards]",
    "description": "Dashboard layout with summary metrics and navigation."
  },
  "created_at": "2026-06-06T10:34:00Z"
}`,
		ErrorResponses: generationErrors("wireframe"),
	},
	{
		Method:              "GET",
		Path:                "/api/history",
		Summary:             "List generation history",
		Description:         "Returns recent generated content history items.",
		Tags:                []string{"History"},
		ResponseDescription: "Recent generation history.",
		ResponseExample: `{
  "items": [
    {
      "id": "b6a0ce5a-a43f-414f-83b4-8a908739bcfb",
      "transcript_id": "2f4fb8e8-8b51-47d0-b699-3fc2d64c6f2d",
      "transcript_title": "Q1 Planning Meeting",
      "type": "prd",
      "created_at": "2026-06-06T10:31:00Z"
    }
  ]
}`,
		ErrorResponses: map[string]string{
			"500": "Failed to get history.",
		},
	},
	{
		Method:              "GET",
		Path:                "/api/history/{id}",
		Summary:             "Get generation by ID",
		Description:         "Returns generated content by generation ID.",
		Tags:                []string{"History"},
		ResponseDescription: "Generated content details.",
		ResponseExample: `{
  "id": "b6a0ce5a-a43f-414f-83b4-8a908739bcfb",
  "transcript_id": "2f4fb8e8-8b51-47d0-b699-3fc2d64c6f2d",
  "type": "prd",
  "content": {
    "content": "# Product Requirements Document\n\n## Overview\nBuild a dashboard for product managers."
  },
  "created_at": "2026-06-06T10:31:00Z"
}`,
		ErrorResponses: map[string]string{
			"400": "Invalid generation ID.",
			"404": "Generation not found.",
		},
	},
}

func generationErrors(kind string) map[string]string {
	return map[string]string{
		"400": "Invalid request body.",
		"404": "Transcript not found.",
		"500": fmt.Sprintf("Failed to generate %s.", kind),
	}
}

func Endpoints() []Endpoint {
	out := make([]Endpoint, len(endpoints))
	copy(out, endpoints)
	return out
}

func OpenAPIJSON() ([]byte, error) {
	return json.MarshalIndent(OpenAPISpec(), "", "  ")
}

func OpenAPISpec() map[string]interface{} {
	paths := map[string]interface{}{}
	for _, endpoint := range endpoints {
		pathItem, _ := paths[endpoint.Path].(map[string]interface{})
		if pathItem == nil {
			pathItem = map[string]interface{}{}
			paths[endpoint.Path] = pathItem
		}

		operation := map[string]interface{}{
			"summary":     endpoint.Summary,
			"description": endpoint.Description,
			"tags":        endpoint.Tags,
			"responses":   responsesFor(endpoint),
		}

		if strings.Contains(endpoint.Path, "{id}") {
			operation["parameters"] = []map[string]interface{}{
				{
					"name":        "id",
					"in":          "path",
					"required":    true,
					"description": "UUID identifier.",
					"schema": map[string]interface{}{
						"type":   "string",
						"format": "uuid",
					},
				},
			}
		}

		if endpoint.RequestContentType != "" {
			operation["requestBody"] = requestBodyFor(endpoint)
		}

		pathItem[strings.ToLower(endpoint.Method)] = operation
	}

	return map[string]interface{}{
		"openapi": "3.0.3",
		"info": map[string]interface{}{
			"title":       "AI PRD Generator API",
			"version":     "1.0.0",
			"description": "API contracts for the AI Product Requirement & UI Generator backend.",
		},
		"servers": []map[string]interface{}{
			{"url": "http://localhost:8080", "description": "Local development"},
		},
		"tags": []map[string]string{
			{"name": "Health", "description": "Service health endpoints"},
			{"name": "Transcripts", "description": "Transcript upload and retrieval"},
			{"name": "Generation", "description": "AI generation endpoints"},
			{"name": "History", "description": "Generation history endpoints"},
		},
		"paths":      paths,
		"components": components(),
	}
}

func requestBodyFor(endpoint Endpoint) map[string]interface{} {
	if endpoint.RequestContentType == "multipart/form-data" {
		return map[string]interface{}{
			"required":    true,
			"description": endpoint.RequestDescription,
			"content": map[string]interface{}{
				"multipart/form-data": map[string]interface{}{
					"schema": map[string]interface{}{
						"type":     "object",
						"required": []string{"content"},
						"properties": map[string]interface{}{
							"title": map[string]interface{}{
								"type":        "string",
								"description": "Optional transcript title.",
								"example":     "Q1 Planning Meeting",
							},
							"content": map[string]interface{}{
								"type":        "string",
								"format":      "binary",
								"description": "Required uploaded transcript file. Permitted file extensions: `.txt`, `.vtt`.",
							},
						},
					},
				},
			},
		}
	}

	return map[string]interface{}{
		"required":    true,
		"description": endpoint.RequestDescription,
		"content": map[string]interface{}{
			"application/json": map[string]interface{}{
				"schema":  map[string]interface{}{"$ref": "#/components/schemas/GenerateRequest"},
				"example": mustJSON(endpoint.RequestExample),
			},
		},
	}
}

func responsesFor(endpoint Endpoint) map[string]interface{} {
	responses := map[string]interface{}{
		"200": map[string]interface{}{
			"description": endpoint.ResponseDescription,
			"content": map[string]interface{}{
				"application/json": map[string]interface{}{
					"schema":  responseSchemaFor(endpoint),
					"example": mustJSON(endpoint.ResponseExample),
				},
			},
		},
	}

	if endpoint.Method == "POST" && endpoint.Path == "/api/transcript" {
		responses["201"] = responses["200"]
		delete(responses, "200")
	}

	for status, description := range endpoint.ErrorResponses {
		responses[status] = map[string]interface{}{
			"description": description,
			"content": map[string]interface{}{
				"application/json": map[string]interface{}{
					"schema": map[string]interface{}{"$ref": "#/components/schemas/ErrorResponse"},
					"example": map[string]interface{}{
						"error": description,
					},
				},
			},
		}
	}

	return responses
}

func responseSchemaFor(endpoint Endpoint) map[string]interface{} {
	switch {
	case endpoint.Path == "/health":
		return map[string]interface{}{"$ref": "#/components/schemas/HealthResponse"}
	case strings.HasPrefix(endpoint.Path, "/api/transcript"):
		return map[string]interface{}{"$ref": "#/components/schemas/Transcript"}
	case strings.HasPrefix(endpoint.Path, "/api/generate") || endpoint.Path == "/api/history/{id}":
		return map[string]interface{}{"$ref": "#/components/schemas/Generation"}
	case endpoint.Path == "/api/history":
		return map[string]interface{}{"$ref": "#/components/schemas/HistoryResponse"}
	default:
		return map[string]interface{}{"type": "object"}
	}
}

func components() map[string]interface{} {
	return map[string]interface{}{
		"schemas": map[string]interface{}{
			"HealthResponse": map[string]interface{}{
				"type":     "object",
				"required": []string{"status", "service"},
				"properties": map[string]interface{}{
					"status":  map[string]string{"type": "string", "example": "healthy"},
					"service": map[string]string{"type": "string", "example": "prd-generator-api"},
				},
			},
			"Transcript": map[string]interface{}{
				"type":     "object",
				"required": []string{"id", "title", "content", "created_at"},
				"properties": map[string]interface{}{
					"id":         map[string]string{"type": "string", "format": "uuid"},
					"title":      map[string]string{"type": "string"},
					"content":    map[string]string{"type": "string", "description": "Uploaded transcript file contents converted to a text string."},
					"created_at": map[string]string{"type": "string", "format": "date-time"},
				},
			},
			"GenerateRequest": map[string]interface{}{
				"type":     "object",
				"required": []string{"transcript_id"},
				"properties": map[string]interface{}{
					"transcript_id": map[string]string{"type": "string", "format": "uuid"},
				},
			},
			"Generation": map[string]interface{}{
				"type":     "object",
				"required": []string{"id", "transcript_id", "type", "content", "created_at"},
				"properties": map[string]interface{}{
					"id":            map[string]string{"type": "string", "format": "uuid"},
					"transcript_id": map[string]string{"type": "string", "format": "uuid"},
					"type":          map[string]interface{}{"type": "string", "enum": []string{"prd", "stories", "requirements", "wireframe"}},
					"content":       map[string]interface{}{"type": "object", "additionalProperties": true},
					"created_at":    map[string]string{"type": "string", "format": "date-time"},
				},
			},
			"HistoryItem": map[string]interface{}{
				"type":     "object",
				"required": []string{"id", "transcript_id", "transcript_title", "type", "created_at"},
				"properties": map[string]interface{}{
					"id":               map[string]string{"type": "string", "format": "uuid"},
					"transcript_id":    map[string]string{"type": "string", "format": "uuid"},
					"transcript_title": map[string]string{"type": "string"},
					"type":             map[string]interface{}{"type": "string", "enum": []string{"prd", "stories", "requirements", "wireframe"}},
					"created_at":       map[string]string{"type": "string", "format": "date-time"},
				},
			},
			"HistoryResponse": map[string]interface{}{
				"type":     "object",
				"required": []string{"items"},
				"properties": map[string]interface{}{
					"items": map[string]interface{}{
						"type":  "array",
						"items": map[string]interface{}{"$ref": "#/components/schemas/HistoryItem"},
					},
				},
			},
			"ErrorResponse": map[string]interface{}{
				"type":     "object",
				"required": []string{"error"},
				"properties": map[string]interface{}{
					"error": map[string]string{"type": "string"},
				},
			},
		},
	}
}

func mustJSON(raw string) interface{} {
	var value interface{}
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		return raw
	}
	return value
}

func Markdown() string {
	var b strings.Builder
	b.WriteString("# AI PRD Generator API Contracts\n\n")
	b.WriteString("> Generated from `backend/internal/apidocs/apidocs.go`. Do not edit manually. Run `go generate ./...` from `backend` after API contract changes.\n\n")
	b.WriteString("Swagger/OpenAPI JSON is generated at `docs/openapi.json` and served by the backend at `/docs/openapi.json`. Swagger UI is available at `/docs`.\n\n")
	b.WriteString("## Base URL\n\n")
	b.WriteString("```text\nhttp://localhost:8080\n```\n\n")
	b.WriteString("## Endpoints\n\n")

	for _, endpoint := range endpoints {
		b.WriteString(fmt.Sprintf("### `%s %s`\n\n", endpoint.Method, endpoint.Path))
		b.WriteString(fmt.Sprintf("**Summary:** %s\n\n", endpoint.Summary))
		b.WriteString(endpoint.Description + "\n\n")

		if endpoint.RequestContentType == "" {
			b.WriteString("#### Request\n\nNo request body.\n\n")
		} else {
			b.WriteString("#### Request\n\n")
			b.WriteString(fmt.Sprintf("Content-Type: `%s`\n\n", endpoint.RequestContentType))
			b.WriteString(endpoint.RequestDescription + "\n\n")
			b.WriteString("Example:\n\n")
			language := "json"
			if endpoint.RequestContentType == "multipart/form-data" {
				language = "bash"
			}
			b.WriteString(fmt.Sprintf("```%s\n%s\n```\n\n", language, endpoint.RequestExample))
		}

		b.WriteString("#### Success Response\n\n")
		status := "200 OK"
		if endpoint.Method == "POST" && endpoint.Path == "/api/transcript" {
			status = "201 Created"
		}
		b.WriteString(fmt.Sprintf("Status: `%s`\n\n", status))
		b.WriteString(endpoint.ResponseDescription + "\n\n")
		b.WriteString("```json\n" + endpoint.ResponseExample + "\n```\n\n")

		if len(endpoint.ErrorResponses) > 0 {
			b.WriteString("#### Error Responses\n\n")
			statuses := make([]string, 0, len(endpoint.ErrorResponses))
			for status := range endpoint.ErrorResponses {
				statuses = append(statuses, status)
			}
			sort.Strings(statuses)
			for _, status := range statuses {
				b.WriteString(fmt.Sprintf("- `%s`: %s\n", status, endpoint.ErrorResponses[status]))
			}
			b.WriteString("\n")
		}
	}

	return b.String()
}
