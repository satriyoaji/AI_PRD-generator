# AI PRD Generator — Frontend

Next.js 16 frontend for the AI Product Requirement & UI Generator. Lets users upload meeting transcript files and generate PRDs, User Stories, Functional Requirements, and UI wireframes through a clean tab-based interface.

## Tech Stack

- **Framework:** Next.js 16 (App Router)
- **Language:** TypeScript
- **Styling:** Tailwind CSS v4 + [shadcn/ui](https://ui.shadcn.com)
- **Fonts:** [Geist](https://vercel.com/font) via `next/font`
- **Icons:** [Lucide React](https://lucide.dev)
- **Notifications:** [Sonner](https://sonner.emilkowal.ski)

## Getting Started

### 1. Install dependencies

```bash
npm install
```

### 2. Configure environment

Create a `.env.local` file in this directory (see [Environment Variables](#environment-variables) below):

```bash
cp .env.local.example .env.local   # if the example file exists, otherwise:
echo "NEXT_PUBLIC_API_URL=http://localhost:8080" > .env.local
```

### 3. Start the backend

The frontend is a thin client — it calls the Go backend for all AI work and data persistence.
Make sure the backend is running before using the app:

```bash
# From the repo root:
cd ../backend && go run ./cmd/server
```

### 4. Run the development server

```bash
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) in your browser.

### Other available scripts

| Command | What it does |
|---|---|
| `npm run dev` | Start development server with hot-reload |
| `npm run build` | Production build |
| `npm run start` | Serve the production build |
| `npm run lint` | Run ESLint |

## Environment Variables

All environment variables used by this app are listed below.
Create a `.env.local` file in `frontend/` and set the values for your environment.
This file is **never committed** to version control (it is in `.gitignore`).

### `NEXT_PUBLIC_API_URL`

| | |
|---|---|
| **Type** | `string` (URL) |
| **Required** | No |
| **Default** | `http://localhost:8080` |
| **Example** | `https://my-backend.railway.app` |

The base URL of the running Go backend. Every API call made by the frontend is prefixed with this value.

- For **local development**: leave it unset or set to `http://localhost:8080`
- For **production / staging**: set it to your deployed backend URL (e.g. Railway, Render)
- The `NEXT_PUBLIC_` prefix is required so Next.js exposes the value to browser-side code

**`.env.local` example:**

```env
# Base URL of the Go backend API.
# Change this when deploying to staging or production.
NEXT_PUBLIC_API_URL=http://localhost:8080
```

> **Note:** After changing `.env.local` you must restart the dev server (`npm run dev`) for
> the new value to take effect. Next.js bakes `NEXT_PUBLIC_*` variables into the client bundle
> at build time.

## Project Structure

```
frontend/
├── app/
│   ├── layout.tsx          # Root layout (fonts, metadata)
│   └── page.tsx            # Main page — wires TranscriptInput + OutputViewer + HistoryPanel
├── components/
│   ├── transcript-input.tsx  # Drag-and-drop / click-to-browse file upload (.txt, .vtt)
│   ├── output-viewer.tsx     # Tab-based viewer for PRD / Stories / Requirements / Wireframe
│   ├── history-panel.tsx     # Recent generation history list
│   └── ui/                   # shadcn/ui primitives (Button, Card, Badge, …)
├── lib/
│   └── api.ts              # Typed API client (fetch wrappers for every backend endpoint)
├── .env.local              # ← your local secrets (git-ignored)
├── .env.local.example      # Safe template to commit and share
└── package.json
```

## Transcript Upload

`POST /api/transcript` accepts **uploaded files only** via `multipart/form-data`.

The file picker in `TranscriptInput` enforces allowed formats on the client side and the
backend validates the extension server-side before accepting the file.

| Field | Type | Required | Description |
|---|---|---|---|
| `content` | file | ✅ | The transcript file to upload |
| `title` | text | ❌ | Human-readable title; defaults to the filename stem |

**Accepted file formats:**

| Extension | Format |
|---|---|
| `.txt` | Plain text transcript |
| `.vtt` | WebVTT (Web Video Text Tracks) — subtitles / captions exported from Zoom, Teams, etc. |

## API Client (`lib/api.ts`)

All backend calls are centralised in `lib/api.ts`. Key methods:

```ts
// Upload a .txt or .vtt file and create a transcript
api.createTranscript(file: File, title?: string): Promise<Transcript>

// Generate documentation from a transcript ID
api.generatePRD(transcriptId: string): Promise<Generation>
api.generateStories(transcriptId: string): Promise<Generation>
api.generateRequirements(transcriptId: string): Promise<Generation>
api.generateWireframe(transcriptId: string): Promise<Generation>

// Retrieve data
api.getTranscript(id: string): Promise<Transcript>
api.getHistory(): Promise<{ items: HistoryItem[] }>
api.getGeneration(id: string): Promise<Generation>
```

For the full interactive API reference, visit the backend Swagger UI at
[http://localhost:8080/docs](http://localhost:8080/docs) while the backend is running.

## Deployment (Vercel)

1. Push your code to GitHub
2. Import the repo into [Vercel](https://vercel.com/new)
3. Set the **Root Directory** to `frontend`
4. Add the environment variable in the Vercel dashboard:
   - `NEXT_PUBLIC_API_URL` → your deployed backend URL
5. Deploy

> The backend must be deployed separately (Railway / Render) and its public URL set as
> `NEXT_PUBLIC_API_URL` before the frontend will function in production.
