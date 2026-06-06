"use client";

import { useState } from "react";
import { Toaster, toast } from "sonner";
import { TranscriptInput } from "@/components/transcript-input";
import { OutputViewer } from "@/components/output-viewer";
import { HistoryPanel } from "@/components/history-panel";
import { Transcript } from "@/lib/api";

export default function Home() {
  const [transcript, setTranscript] = useState<Transcript | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [refreshKey, setRefreshKey] = useState(0);

  const handleTranscriptCreated = (newTranscript: Transcript) => {
    setTranscript(newTranscript);
    toast.success("Transcript created successfully!");
    setRefreshKey((k) => k + 1);
  };

  return (
    <div className="min-h-screen bg-gradient-to-b from-slate-50 to-slate-100 dark:from-slate-950 dark:to-slate-900">
      {/* Header */}
      <header className="border-b bg-white/80 backdrop-blur-sm dark:bg-slate-950/80">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-2xl font-bold tracking-tight">
                AI PRD Generator
              </h1>
              <p className="text-sm text-muted-foreground">
                Transform meeting transcripts into product documentation
              </p>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="container mx-auto px-4 py-8">
        <div className="grid gap-6 lg:grid-cols-2">
          {/* Left Column - Input & History */}
          <div className="space-y-6">
            <TranscriptInput
              onTranscriptCreated={handleTranscriptCreated}
              isLoading={isLoading}
              setIsLoading={setIsLoading}
            />
            <HistoryPanel
              onSelectTranscript={setTranscript}
              refreshKey={refreshKey}
            />
          </div>

          {/* Right Column - Output */}
          <div className="lg:sticky lg:top-4 lg:h-[calc(100vh-8rem)]">
            <OutputViewer
              transcriptId={transcript?.id || null}
              refreshKey={refreshKey}
            />
          </div>
        </div>
      </main>

      {/* Footer */}
      <footer className="border-t bg-white/50 dark:bg-slate-950/50">
        <div className="container mx-auto px-4 py-4">
          <p className="text-center text-sm text-muted-foreground">
            Built with Go + Next.js | AI-powered Product Documentation Generator
          </p>
        </div>
      </footer>

      <Toaster position="top-right" />
    </div>
  );
}
