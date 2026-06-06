# AI PRD Generator - Technical Assessment Project

This is a full-stack application built for the AI Product Requirement & UI Generator technical assessment.

## Project Structure

```
ai_prd-generator-test/
├── backend/              # Go backend with Fiber
│   ├── cmd/server/       # Main entry point
│   ├── internal/
│   │   ├── config/       # Configuration loading
│   │   ├── db/           # PostgreSQL connection
│   │   ├── handlers/     # HTTP handlers
│   │   ├── models/       # Data models
│   │   └── services/     # AI service layer
│   ├── .env.example
│   ├── Dockerfile
│   └── go.mod
├── frontend/            # Next.js frontend
│   ├── app/              # Next.js app router
│   ├── components/       # React components
│   ├── lib/              # Utilities and API client
│   ├── .env.local
│   └── package.json
└── README.md
```

## Tech Stack

- **Backend:** Go 1.23+, Fiber v2, pgx v4, PostgreSQL
- **Frontend:** Next.js 16, TypeScript, Tailwind CSS, shadcn/ui
- **AI:** Anthropic Claude (PRD/Stories/Requirements), OpenAI GPT-4o (Wireframes)
- **Database:** PostgreSQL (Neon recommended)

## Quick Start

### Backend
```bash
cd backend
cp .env.example .env
# Edit .env with your API keys
go run ./cmd/server
```

### Frontend
```bash
cd frontend
npm install
npm run dev
```

## Key Files

- `backend/cmd/server/main.go` - Server entry point and route setup
- `backend/internal/services/ai_service.go` - AI generation logic
- `frontend/app/page.tsx` - Main application page
- `frontend/components/output-viewer.tsx` - Tab-based output viewer
- `frontend/lib/api.ts` - API client

## Deployment

- **Backend:** Deploy to Railway/Render with DATABASE_URL and API keys
- **Frontend:** Deploy to Vercel with NEXT_PUBLIC_API_URL pointing to backend