const API_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export interface Transcript {
  id: string;
  title: string;
  content: string;
  created_at: string;
}

export interface UserStory {
  id: string;
  title: string;
  description: string;
  acceptance_criteria: string[];
}

export interface Generation {
  id: string;
  transcript_id: string;
  type: string;
  content: any;
  created_at: string;
}

export interface HistoryItem {
  id: string;
  transcript_id: string;
  transcript_title: string;
  type: string;
  created_at: string;
}

export interface PRDContent {
  content: string;
}

export interface RequirementsContent {
  content: string;
}

export interface WireframeContent {
  mermaid_code: string;
  description: string;
}

export interface StoriesContent {
  stories: UserStory[];
}

async function fetchAPI<T>(
  endpoint: string,
  options?: RequestInit,
): Promise<T> {
  const response = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...options?.headers,
    },
  });

  if (!response.ok) {
    const error = await response
      .json()
      .catch(() => ({ error: "Request failed" }));
    throw new Error(error.error || "Request failed");
  }

  return response.json();
}

/**
 * Upload a transcript file (.txt or .vtt) via multipart/form-data.
 * The browser sets the correct Content-Type boundary automatically when
 * FormData is passed as the body, so we must NOT set Content-Type manually.
 */
async function uploadTranscript(
  file: File,
  title?: string,
): Promise<Transcript> {
  const form = new FormData();
  form.append("content", file);
  if (title) form.append("title", title);

  const response = await fetch(`${API_BASE}/api/transcript`, {
    method: "POST",
    body: form,
    // intentionally no Content-Type header — browser fills in multipart boundary
  });

  if (!response.ok) {
    const error = await response
      .json()
      .catch(() => ({ error: "Upload failed" }));
    throw new Error(error.error || "Upload failed");
  }

  return response.json();
}

export const api = {
  // Transcript endpoints — POST uses multipart file upload
  createTranscript: (file: File, title?: string) =>
    uploadTranscript(file, title),

  getTranscript: (id: string) => fetchAPI<Transcript>(`/api/transcript/${id}`),

  // Generate endpoints
  generatePRD: (transcriptId: string) =>
    fetchAPI<Generation>("/api/generate/prd", {
      method: "POST",
      body: JSON.stringify({ transcript_id: transcriptId }),
    }),

  generateStories: (transcriptId: string) =>
    fetchAPI<Generation>("/api/generate/stories", {
      method: "POST",
      body: JSON.stringify({ transcript_id: transcriptId }),
    }),

  generateRequirements: (transcriptId: string) =>
    fetchAPI<Generation>("/api/generate/requirements", {
      method: "POST",
      body: JSON.stringify({ transcript_id: transcriptId }),
    }),

  generateWireframe: (transcriptId: string) =>
    fetchAPI<Generation>("/api/generate/wireframe", {
      method: "POST",
      body: JSON.stringify({ transcript_id: transcriptId }),
    }),

  // History endpoints
  getHistory: () => fetchAPI<{ items: HistoryItem[] }>("/api/history"),

  getGeneration: (id: string) => fetchAPI<Generation>(`/api/history/${id}`),

  // Health check
  health: () => fetchAPI<{ status: string; service: string }>("/health"),
};
