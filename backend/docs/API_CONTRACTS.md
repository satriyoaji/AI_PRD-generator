# AI PRD Generator API Contracts

> Generated from `backend/internal/apidocs/apidocs.go`. Do not edit manually. Run `go generate ./...` from `backend` after API contract changes.

Swagger/OpenAPI JSON is generated at `docs/openapi.json` and served by the backend at `/docs/openapi.json`. Swagger UI is available at `/docs`.

## Base URL

```text
http://localhost:8080
```

## Endpoints

### `GET /health`

**Summary:** Health check

Returns the API health status.

#### Request

No request body.

#### Success Response

Status: `200 OK`

Service health status.

```json
{
  "status": "healthy",
  "service": "prd-generator-api"
}
```

### `POST /api/transcript`

**Summary:** Create transcript from uploaded file

Creates a transcript by uploading a plain text `.txt` file or WebVTT `.vtt` file. The uploaded file must be sent as multipart form-data using the `content` file field. The file body is read and stored as a text string.

#### Request

Content-Type: `multipart/form-data`

Multipart form with optional `title` text field and required `content` file field. Only `.txt` and `.vtt` file extensions are permitted.

Example:

```bash
curl -X POST http://localhost:8080/api/transcript \
  -F "title=Q1 Planning Meeting" \
  -F "content=@meeting.vtt"
```

#### Success Response

Status: `201 Created`

Created transcript. `content` contains the uploaded file contents converted to a string.

```json
{
  "id": "2f4fb8e8-8b51-47d0-b699-3fc2d64c6f2d",
  "title": "Q1 Planning Meeting",
  "content": "WEBVTT\n\n00:00:00.000 --> 00:00:05.000\nWe need a dashboard for product managers.",
  "created_at": "2026-06-06T10:30:00Z"
}
```

#### Error Responses

- `400`: Missing file, empty file, or malformed multipart request.
- `415`: Uploaded file extension is not `.txt` or `.vtt`.
- `500`: Failed to read or save transcript.

### `GET /api/transcript/{id}`

**Summary:** Get transcript by ID

Returns a previously created transcript.

#### Request

No request body.

#### Success Response

Status: `200 OK`

Transcript details.

```json
{
  "id": "2f4fb8e8-8b51-47d0-b699-3fc2d64c6f2d",
  "title": "Q1 Planning Meeting",
  "content": "We need a dashboard for product managers.",
  "created_at": "2026-06-06T10:30:00Z"
}
```

#### Error Responses

- `400`: Invalid transcript ID.
- `404`: Transcript not found.

### `POST /api/generate/prd`

**Summary:** Generate PRD

Generates a Product Requirements Document from a transcript.

#### Request

Content-Type: `application/json`

JSON body containing the transcript ID.

Example:

```json
{
  "transcript_id": "2f4fb8e8-8b51-47d0-b699-3fc2d64c6f2d"
}
```

#### Success Response

Status: `200 OK`

Generated PRD content.

```json
{
  "id": "b6a0ce5a-a43f-414f-83b4-8a908739bcfb",
  "transcript_id": "2f4fb8e8-8b51-47d0-b699-3fc2d64c6f2d",
  "type": "prd",
  "content": {
    "content": "# Product Requirements Document\n\n## Overview\nBuild a dashboard for product managers."
  },
  "created_at": "2026-06-06T10:31:00Z"
}
```

#### Error Responses

- `400`: Invalid request body.
- `404`: Transcript not found.
- `500`: Failed to generate PRD.

### `POST /api/generate/stories`

**Summary:** Generate user stories

Generates user stories from a transcript.

#### Request

Content-Type: `application/json`

JSON body containing the transcript ID.

Example:

```json
{
  "transcript_id": "2f4fb8e8-8b51-47d0-b699-3fc2d64c6f2d"
}
```

#### Success Response

Status: `200 OK`

Generated user stories.

```json
{
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
}
```

#### Error Responses

- `400`: Invalid request body.
- `404`: Transcript not found.
- `500`: Failed to generate user stories.

### `POST /api/generate/requirements`

**Summary:** Generate functional requirements

Generates functional requirements from a transcript.

#### Request

Content-Type: `application/json`

JSON body containing the transcript ID.

Example:

```json
{
  "transcript_id": "2f4fb8e8-8b51-47d0-b699-3fc2d64c6f2d"
}
```

#### Success Response

Status: `200 OK`

Generated functional requirements.

```json
{
  "id": "a260d2d4-d2aa-4e55-8765-51d1530748af",
  "transcript_id": "2f4fb8e8-8b51-47d0-b699-3fc2d64c6f2d",
  "type": "requirements",
  "content": {
    "content": "# Functional Requirements\n\n1. The system shall display dashboard metrics."
  },
  "created_at": "2026-06-06T10:33:00Z"
}
```

#### Error Responses

- `400`: Invalid request body.
- `404`: Transcript not found.
- `500`: Failed to generate requirements.

### `POST /api/generate/wireframe`

**Summary:** Generate wireframe

Generates UI/UX wireframe documentation from a transcript.

#### Request

Content-Type: `application/json`

JSON body containing the transcript ID.

Example:

```json
{
  "transcript_id": "2f4fb8e8-8b51-47d0-b699-3fc2d64c6f2d"
}
```

#### Success Response

Status: `200 OK`

Generated wireframe content.

```json
{
  "id": "6957633d-7716-46b7-8bd3-3d0d87e8a318",
  "transcript_id": "2f4fb8e8-8b51-47d0-b699-3fc2d64c6f2d",
  "type": "wireframe",
  "content": {
    "mermaid_code": "graph TD\n  A[Dashboard] --> B[Metrics Cards]",
    "description": "Dashboard layout with summary metrics and navigation."
  },
  "created_at": "2026-06-06T10:34:00Z"
}
```

#### Error Responses

- `400`: Invalid request body.
- `404`: Transcript not found.
- `500`: Failed to generate wireframe.

### `GET /api/history`

**Summary:** List generation history

Returns recent generated content history items.

#### Request

No request body.

#### Success Response

Status: `200 OK`

Recent generation history.

```json
{
  "items": [
    {
      "id": "b6a0ce5a-a43f-414f-83b4-8a908739bcfb",
      "transcript_id": "2f4fb8e8-8b51-47d0-b699-3fc2d64c6f2d",
      "transcript_title": "Q1 Planning Meeting",
      "type": "prd",
      "created_at": "2026-06-06T10:31:00Z"
    }
  ]
}
```

#### Error Responses

- `500`: Failed to get history.

### `GET /api/history/{id}`

**Summary:** Get generation by ID

Returns generated content by generation ID.

#### Request

No request body.

#### Success Response

Status: `200 OK`

Generated content details.

```json
{
  "id": "b6a0ce5a-a43f-414f-83b4-8a908739bcfb",
  "transcript_id": "2f4fb8e8-8b51-47d0-b699-3fc2d64c6f2d",
  "type": "prd",
  "content": {
    "content": "# Product Requirements Document\n\n## Overview\nBuild a dashboard for product managers."
  },
  "created_at": "2026-06-06T10:31:00Z"
}
```

#### Error Responses

- `400`: Invalid generation ID.
- `404`: Generation not found.

