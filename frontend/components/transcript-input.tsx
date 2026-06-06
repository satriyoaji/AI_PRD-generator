"use client";

import { useRef, useState } from "react";
import { Upload, FileText, X, CheckCircle2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { api, Transcript } from "@/lib/api";

interface TranscriptInputProps {
  onTranscriptCreated: (transcript: Transcript) => void;
  isLoading: boolean;
  setIsLoading: (loading: boolean) => void;
}

const ACCEPTED_EXTENSIONS = [".txt", ".vtt"];
const ACCEPTED_MIME = ["text/plain", "text/vtt"];

function isAllowedFile(file: File): boolean {
  const name = file.name.toLowerCase();
  return ACCEPTED_EXTENSIONS.some((ext) => name.endsWith(ext));
}

export function TranscriptInput({
  onTranscriptCreated,
  isLoading,
  setIsLoading,
}: TranscriptInputProps) {
  const [title, setTitle] = useState("");
  const [file, setFile] = useState<File | null>(null);
  const [preview, setPreview] = useState<string>("");
  const [error, setError] = useState<string | null>(null);
  const [isDragging, setIsDragging] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  // ── helpers ───────────────────────────────────────────────────────────────

  const applyFile = (picked: File) => {
    if (!isAllowedFile(picked)) {
      setError(
        `Unsupported file type "${picked.name}". Only .txt and .vtt files are accepted.`,
      );
      return;
    }
    setError(null);
    setFile(picked);

    // auto-populate title from filename stem
    if (!title) {
      const stem = picked.name.replace(/\.[^/.]+$/, "");
      setTitle(stem);
    }

    // generate preview
    const reader = new FileReader();
    reader.onload = (e) => {
      const text = (e.target?.result as string) ?? "";
      setPreview(text.slice(0, 400) + (text.length > 400 ? "…" : ""));
    };
    reader.readAsText(picked);
  };

  const clearFile = () => {
    setFile(null);
    setPreview("");
    setError(null);
    if (inputRef.current) inputRef.current.value = "";
  };

  // ── event handlers ────────────────────────────────────────────────────────

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const picked = e.target.files?.[0];
    if (picked) applyFile(picked);
  };

  const handleDrop = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsDragging(false);
    const picked = e.dataTransfer.files?.[0];
    if (picked) applyFile(picked);
  };

  const handleDragOver = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = () => setIsDragging(false);

  const handleSubmit = async () => {
    if (!file) {
      setError("Please upload a .txt or .vtt transcript file.");
      return;
    }

    setError(null);
    setIsLoading(true);

    try {
      const transcript = await api.createTranscript(file, title || undefined);
      onTranscriptCreated(transcript);
      // reset form
      clearFile();
      setTitle("");
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to create transcript",
      );
    } finally {
      setIsLoading(false);
    }
  };

  // ── render ────────────────────────────────────────────────────────────────

  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <FileText className="h-5 w-5" />
          Meeting Transcript
        </CardTitle>
        <CardDescription>
          Upload a <code>.txt</code> or <code>.vtt</code> transcript file to
          generate product documentation
        </CardDescription>
      </CardHeader>

      <CardContent className="space-y-4">
        {/* Title */}
        <div className="space-y-2">
          <label htmlFor="title" className="text-sm font-medium">
            Title{" "}
            <span className="text-muted-foreground font-normal">
              (optional)
            </span>
          </label>
          <Input
            id="title"
            placeholder="e.g., Q1 Planning Meeting"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
          />
        </div>

        {/* Drop zone */}
        <div
          onDrop={handleDrop}
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
          onClick={() => !file && inputRef.current?.click()}
          className={[
            "relative flex flex-col items-center justify-center gap-3 rounded-lg border-2 border-dashed p-8 text-center transition-colors",
            isDragging
              ? "border-primary bg-primary/5"
              : "border-muted-foreground/25 hover:border-primary/50 hover:bg-muted/30",
            !file ? "cursor-pointer" : "",
          ].join(" ")}
        >
          {file ? (
            // ── file selected state ─────────────────────────────────────────
            <>
              <CheckCircle2 className="h-8 w-8 text-green-500 shrink-0" />
              <div className="flex flex-col items-center gap-1">
                <p className="text-sm font-medium">{file.name}</p>
                <p className="text-xs text-muted-foreground">
                  {(file.size / 1024).toFixed(1)} KB
                </p>
              </div>

              {preview && (
                <pre className="w-full max-h-28 overflow-hidden rounded-md bg-muted/60 p-3 text-left text-xs font-mono text-muted-foreground whitespace-pre-wrap">
                  {preview}
                </pre>
              )}

              <Button
                type="button"
                variant="ghost"
                size="sm"
                className="absolute top-2 right-2 h-7 w-7 p-0 text-muted-foreground hover:text-destructive"
                onClick={(e) => {
                  e.stopPropagation();
                  clearFile();
                }}
              >
                <X className="h-4 w-4" />
                <span className="sr-only">Remove file</span>
              </Button>
            </>
          ) : (
            // ── empty / prompt state ────────────────────────────────────────
            <>
              <Upload className="h-8 w-8 text-muted-foreground" />
              <div className="space-y-1">
                <p className="text-sm font-medium">
                  Drop your transcript here, or{" "}
                  <span className="text-primary underline underline-offset-2">
                    browse
                  </span>
                </p>
                <p className="text-xs text-muted-foreground">
                  Accepted formats:
                </p>
                <div className="flex justify-center gap-2">
                  <Badge variant="secondary">.txt — Plain text</Badge>
                  <Badge variant="secondary">.vtt — WebVTT</Badge>
                </div>
              </div>
            </>
          )}

          {/* hidden native input */}
          <input
            ref={inputRef}
            type="file"
            accept={ACCEPTED_EXTENSIONS.join(",")}
            className="hidden"
            onChange={handleFileChange}
          />
        </div>

        {/* Error */}
        {error && <p className="text-sm text-destructive">{error}</p>}

        {/* Submit */}
        <Button
          onClick={handleSubmit}
          disabled={isLoading || !file}
          className="w-full"
        >
          {isLoading ? "Uploading…" : "Upload & Create Transcript"}
        </Button>
      </CardContent>
    </Card>
  );
}
