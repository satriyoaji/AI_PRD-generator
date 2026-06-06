"use client";

import { useEffect, useState, useCallback } from "react";
import { History, RefreshCw } from "lucide-react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { ScrollArea } from "@/components/ui/scroll-area";
import { api, HistoryItem, Transcript } from "@/lib/api";

interface HistoryPanelProps {
  onSelectTranscript: (transcript: Transcript) => void;
  refreshKey: number;
}

export function HistoryPanel({ onSelectTranscript, refreshKey }: HistoryPanelProps) {
  const [history, setHistory] = useState<HistoryItem[]>([]);
  const [loading, setLoading] = useState(false);

  const fetchHistory = useCallback(async () => {
    setLoading(true);
    try {
      const response = await api.getHistory();
      setHistory(response.items || []);
    } catch (error) {
      console.error("Failed to fetch history:", error);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    const loadData = () => {
      fetchHistory();
    };
    loadData();
  }, [refreshKey, fetchHistory]);

  const getTypeLabel = (type: string) => {
    const labels: Record<string, string> = {
      prd: "PRD",
      stories: "Stories",
      requirements: "Reqs",
      wireframe: "Wireframe",
    };
    return labels[type] || type;
  };

  const getTypeColor = (type: string) => {
    const colors: Record<string, string> = {
      prd: "bg-blue-500/10 text-blue-600 border-blue-500/20",
      stories: "bg-green-500/10 text-green-600 border-green-500/20",
      requirements: "bg-amber-500/10 text-amber-600 border-amber-500/20",
      wireframe: "bg-purple-500/10 text-purple-600 border-purple-500/20",
    };
    return colors[type] || "";
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  return (
    <Card className="w-full">
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <CardTitle className="flex items-center gap-2 text-base">
            <History className="h-4 w-4" />
            History
          </CardTitle>
          <Button size="sm" variant="ghost" onClick={fetchHistory} disabled={loading}>
            <RefreshCw className={`h-4 w-4 ${loading ? "animate-spin" : ""}`} />
          </Button>
        </div>
        <CardDescription>Recent generations from your transcripts</CardDescription>
      </CardHeader>

      <CardContent className="p-0">
        <ScrollArea className="h-[300px]">
          {history.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-8 text-center px-4">
              <History className="h-8 w-8 text-muted-foreground mb-2" />
              <p className="text-sm text-muted-foreground">No history yet</p>
              <p className="text-xs text-muted-foreground mt-1">
                Generated content will appear here
              </p>
            </div>
          ) : (
            <div className="space-y-1 p-2">
              {history.map((item) => (
                <button
                  key={item.id}
                  onClick={() => {
                    api.getTranscript(item.transcript_id).then(onSelectTranscript).catch(console.error);
                  }}
                  className="w-full text-left p-3 rounded-lg hover:bg-muted/50 transition-colors"
                >
                  <div className="flex items-start justify-between gap-2">
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium truncate">
                        {item.transcript_title || "Untitled Transcript"}
                      </p>
                      <p className="text-xs text-muted-foreground mt-0.5">
                        {formatDate(item.created_at)}
                      </p>
                    </div>
                    <Badge
                      variant="outline"
                      className={`text-xs shrink-0 ${getTypeColor(item.type)}`}
                    >
                      {getTypeLabel(item.type)}
                    </Badge>
                  </div>
                </button>
              ))}
            </div>
          )}
        </ScrollArea>
      </CardContent>
    </Card>
  );
}