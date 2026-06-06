# AI Product Requirement & UI Generator

A full-stack application that converts meeting transcripts into Product Requirements Documents (PRD), User Stories, Functional Requirements, and UI/UX Wireframes using AI.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        Frontend (Next.js)                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────┐ │
│  │ Transcript │  │   Output    │  │   History   │  │ Export  │ │
│  │   Upload   │  │   Viewer    │  │   Panel     │  │  Panel  │ │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────┘ │
└─────────────────────────────────────────────────────────────────┘
                              │ REST API
┌─────────────────────────────────────────────────────────────────┐
│                      Backend (Go/Fiber)                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────┐ │
│  │   Router   │  │   Service   │  │    AI       │  │   DB    │ │
│  │            │  │   Layer     │  │   Client    │  │  Layer  │ │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────┘ │
└─────────────────────────────────────────────────────────────────┘
                              │
         ┌────────────────────┼────────────────────┐
         ▼                    ▼                    ▼
   ┌──────────┐        ┌──────────┐        ┌──────────┐
   │ Claude   │        │  GPT-4o  │        │PostgreSQL│
   │(Anthropic)│        │ (OpenAI) │        │  (Neon)  │
   └──────────┘        └──────────┘        └──────────┘
```

## Technology Stack

### Backend
- **Language:** Go 1.23+
- **Framework:** Fiber v2 (high-performance HTTP framework)
- **Database:** PostgreSQL via pgx v4
- **AI Providers:** Anthropic Claude API, OpenAI GPT-4o API

### Frontend
- **Framework:** Next.js 16 (App Router)
- **Language:** TypeScript
- **Styling:** Tailwind CSS v4 + shadcn/ui
- **State Management:** React hooks

### Infrastructure
- **Database:** Neon PostgreSQL (free tier)
- **Backend Hosting:** Railway / Render
- **Frontend Hosting:** Vercel

## Features

- **Transcript Upload:** Upload `.txt` or `.vtt` meeting transcript files (drag-and-drop or click-to-browse)
- **PRD Generation:** AI-powered Product Requirements Document creation
- **User Stories:** Generate user stories with acceptance criteria
- **Functional Requirements:** Extract and format functional requirements
- **UI Wireframes:** Generate Mermaid-based wireframe diagrams
- **History Panel:** View and reuse past generations
- **Export:** Copy or download generated content as Markdown

## Setup Instructions

### Prerequisites

- Go 1.23+
- Node.js 18+
- PostgreSQL database (or use Neon cloud)
- Anthropic API key (for Claude)
- OpenAI API key (for GPT-4o)

### Backend Setup

```bash
cd backend

# Copy environment template and fill in your secrets
cp .env.example .env

# Install/tidy Go dependencies
go mod tidy

# Run the development server (port 8080 by default)
go run ./cmd/server
```

### Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Create the local environment file (see frontend/README.md for all variables)
echo "NEXT_PUBLIC_API_URL=http://localhost:8080" > .env.local

# Run the development server (port 3000 by default)
npm run dev
```

### Docker (Optional)

```bash
cd backend

docker build -t prd-generator .
docker run -p 8080:8080 --env-file .env prd-generator
```

## Configuration

### Backend Environment Variables (`backend/.env`)

| Variable | Description | Required |
|---|---|---|
| `PORT` | Server port (default: `8080`) | No |
| `DATABASE_URL` | PostgreSQL connection string | No — falls back to mock data |
| `ANTHROPIC_API_KEY` | Anthropic Claude API key (PRD, Stories, Requirements) | Yes |
| `OPENAI_API_KEY` | OpenAI GPT-4o API key (Wireframes) | Yes |

### Frontend Environment Variables (`frontend/.env.local`)

| Variable | Description | Required |
|---|---|---|
| `NEXT_PUBLIC_API_URL` | Full URL of the running backend, e.g. `http://localhost:8080` | No — defaults to `http://localhost:8080` |

> See [`frontend/README.md`](frontend/README.md) for full frontend environment documentation.

## API Endpoints

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/health` | Health check |
| `POST` | `/api/transcript` | Upload a `.txt` or `.vtt` file and create a transcript |
| `GET` | `/api/transcript/:id` | Retrieve a transcript |
| `POST` | `/api/generate/prd` | Generate PRD from transcript |
| `POST` | `/api/generate/stories` | Generate User Stories from transcript |
| `POST` | `/api/generate/requirements` | Generate Functional Requirements from transcript |
| `POST` | `/api/generate/wireframe` | Generate UI Wireframe from transcript |
| `GET` | `/api/history` | List generation history |
| `GET` | `/api/history/:id` | Get a specific generation |

## API Documentation (Live Swagger + Contracts)

The backend serves interactive API documentation at runtime — always in sync with the code.

| URL | What you get |
|---|---|
| [`http://localhost:8080/docs`](http://localhost:8080/docs) | **Swagger UI** — interactive explorer for every endpoint |
| [`http://localhost:8080/docs/openapi.json`](http://localhost:8080/docs/openapi.json) | **OpenAPI 3.0 JSON** spec (importable into Postman, Insomnia, etc.) |
| [`http://localhost:8080/docs/api-contracts.md`](http://localhost:8080/docs/api-contracts.md) | **Markdown contract reference** served live |

### Static doc files (committed to the repo)

Pre-generated copies live in `backend/docs/` so the contracts are readable without running the server:

```
backend/docs/
├── openapi.json        ← OpenAPI 3.0 spec
└── API_CONTRACTS.md    ← full markdown API reference
```

### Regenerating the static files

Whenever you change an API contract, update the single source of truth:

```
backend/internal/apidocs/apidocs.go
```

Then regenerate the static files:

```bash
cd backend

# Option 1 — go generate (uses the //go:generate directive in cmd/server/main.go)
go generate ./...

# Option 2 — run the generator directly
go run ./cmd/docs
```

Both commands overwrite `backend/docs/openapi.json` and `backend/docs/API_CONTRACTS.md`.
The live server routes (`/docs`, `/docs/openapi.json`, `/docs/api-contracts.md`) always
reflect the latest code without any regeneration step — just restart the server.

## Known Limitations

1. **No Authentication:** No user authentication system
2. **Mock Database Fallback:** When `DATABASE_URL` is not set, uses in-memory mock data that does not persist
3. **API Rate Limits:** Subject to Anthropic and OpenAI rate limits
4. **Transcript File Formats:** Upload accepts `.txt` (plain text) and `.vtt` (WebVTT) only
5. **Wireframe Rendering:** Mermaid diagrams require an external viewer such as [mermaid.live](https://mermaid.live)

## Future Improvements

- [ ] User authentication and authorization
- [ ] Support additional transcript file formats (e.g. `.srt`, `.json`)
- [ ] PDF export functionality
- [ ] Collaborative editing
- [ ] Custom prompt templates
- [ ] Integration with video conferencing tools
- [ ] Mobile responsive improvements and dark-mode toggle
- [ ] Real-time streaming responses
- [ ] Language support for non-English transcripts

## License

MIT License

## Author

Built for the AI Product Requirement & UI Generator Technical Assessment
