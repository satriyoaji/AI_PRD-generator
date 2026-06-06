"use client";

import { useState } from "react";
import {
  FileText,
  ListChecks,
  ListTodo,
  Layout,
  Copy,
  Download,
  CheckCircle,
  Loader2,
} from "lucide-react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { Generation, api, UserStory } from "@/lib/api";

interface OutputViewerProps {
  transcriptId: string | null;
  refreshKey: number;
}

export function OutputViewer({ transcriptId, refreshKey }: OutputViewerProps) {
  const [activeTab, setActiveTab] = useState("prd");
  const [generations, setGenerations] = useState<Record<string, Generation>>(
    {},
  );
  const [loading, setLoading] = useState<Record<string, boolean>>({});
  const [copied, setCopied] = useState<string | null>(null);

  const generate = async (type: string) => {
    if (!transcriptId) return;

    setLoading((prev) => ({ ...prev, [type]: true }));
    try {
      let generation: Generation;
      switch (type) {
        case "prd":
          generation = await api.generatePRD(transcriptId);
          break;
        case "stories":
          generation = await api.generateStories(transcriptId);
          break;
        case "requirements":
          generation = await api.generateRequirements(transcriptId);
          break;
        case "wireframe":
          generation = await api.generateWireframe(transcriptId);
          break;
        default:
          return;
      }
      setGenerations((prev) => ({ ...prev, [type]: generation }));
      setActiveTab(type);
    } catch (error) {
      console.error(`Failed to generate ${type}:`, error);
    } finally {
      setLoading((prev) => ({ ...prev, [type]: false }));
    }
  };

  const copyToClipboard = async (content: string, type: string) => {
    await navigator.clipboard.writeText(content);
    setCopied(type);
    setTimeout(() => setCopied(null), 2000);
  };

  const downloadContent = (content: string, filename: string) => {
    const blob = new Blob([content], { type: "text/markdown" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = filename;
    a.click();
    URL.revokeObjectURL(url);
  };

  const getContentText = (type: string): string => {
    const gen = generations[type];
    if (!gen) return "";

    const content = gen.content;
    if (type === "prd" || type === "requirements") {
      return content.content || "";
    }
    if (type === "stories") {
      return JSON.stringify(content, null, 2);
    }
    if (type === "wireframe") {
      return `## Wireframe Description\n\n${content.description || ""}\n\n## Mermaid Code\n\n\`\`\`mermaid\n${content.mermaid_code || ""}\n\`\`\``;
    }
    return "";
  };

  if (!transcriptId) {
    return (
      <Card className="w-full h-full flex items-center justify-center">
        <CardContent className="text-center py-12">
          <FileText className="h-12 w-12 mx-auto mb-4 text-muted-foreground" />
          <p className="text-muted-foreground">
            Create a transcript to start generating documentation
          </p>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="w-full h-full flex flex-col">
      <CardHeader>
        <CardTitle className="flex items-center justify-between">
          <span>Generated Output</span>
          <Badge variant="outline" className="text-xs">
            {Object.keys(generations).length}/4 Generated
          </Badge>
        </CardTitle>
        <CardDescription>
          Generate and view product documentation from your transcript
        </CardDescription>
      </CardHeader>

      <CardContent className="flex-1 flex flex-col min-h-0">
        <Tabs
          value={activeTab}
          onValueChange={setActiveTab}
          className="flex-1 flex flex-col"
        >
          <TabsList className="grid w-full grid-cols-4">
            <TabsTrigger value="prd" className="flex items-center gap-1">
              <FileText className="h-3 w-3" />
              <span className="hidden sm:inline">PRD</span>
            </TabsTrigger>
            <TabsTrigger value="stories" className="flex items-center gap-1">
              <ListChecks className="h-3 w-3" />
              <span className="hidden sm:inline">Stories</span>
            </TabsTrigger>
            <TabsTrigger
              value="requirements"
              className="flex items-center gap-1"
            >
              <ListTodo className="h-3 w-3" />
              <span className="hidden sm:inline">Reqs</span>
            </TabsTrigger>
            <TabsTrigger value="wireframe" className="flex items-center gap-1">
              <Layout className="h-3 w-3" />
              <span className="hidden sm:inline">UI</span>
            </TabsTrigger>
          </TabsList>

          <div className="flex items-center justify-between mt-4 mb-2">
            <div className="flex gap-2">
              {!generations.prd && (
                <Button
                  size="sm"
                  variant="outline"
                  onClick={() => generate("prd")}
                  disabled={loading.prd}
                >
                  {loading.prd ? (
                    <Loader2 className="h-4 w-4 animate-spin" />
                  ) : (
                    <FileText className="h-4 w-4" />
                  )}
                  <span className="ml-1">Generate PRD</span>
                </Button>
              )}
              {!generations.stories && (
                <Button
                  size="sm"
                  variant="outline"
                  onClick={() => generate("stories")}
                  disabled={loading.stories}
                >
                  {loading.stories ? (
                    <Loader2 className="h-4 w-4 animate-spin" />
                  ) : (
                    <ListChecks className="h-4 w-4" />
                  )}
                  <span className="ml-1">Stories</span>
                </Button>
              )}
              {!generations.requirements && (
                <Button
                  size="sm"
                  variant="outline"
                  onClick={() => generate("requirements")}
                  disabled={loading.requirements}
                >
                  {loading.requirements ? (
                    <Loader2 className="h-4 w-4 animate-spin" />
                  ) : (
                    <ListTodo className="h-4 w-4" />
                  )}
                  <span className="ml-1">Reqs</span>
                </Button>
              )}
              {!generations.wireframe && (
                <Button
                  size="sm"
                  variant="outline"
                  onClick={() => generate("wireframe")}
                  disabled={loading.wireframe}
                >
                  {loading.wireframe ? (
                    <Loader2 className="h-4 w-4 animate-spin" />
                  ) : (
                    <Layout className="h-4 w-4" />
                  )}
                  <span className="ml-1">Wireframe</span>
                </Button>
              )}
            </div>

            {generations[activeTab] && (
              <div className="flex gap-1">
                <Button
                  size="sm"
                  variant="ghost"
                  onClick={() =>
                    copyToClipboard(getContentText(activeTab), activeTab)
                  }
                >
                  {copied === activeTab ? (
                    <CheckCircle className="h-4 w-4 text-green-500" />
                  ) : (
                    <Copy className="h-4 w-4" />
                  )}
                </Button>
                <Button
                  size="sm"
                  variant="ghost"
                  onClick={() =>
                    downloadContent(
                      getContentText(activeTab),
                      `${activeTab}.md`,
                    )
                  }
                >
                  <Download className="h-4 w-4" />
                </Button>
              </div>
            )}
          </div>

          <ScrollArea className="flex-1 min-h-[400px]">
            <TabsContent value="prd" className="mt-0">
              {generations.prd ? (
                <div className="prose prose-sm max-w-none dark:prose-invert">
                  <pre className="whitespace-pre-wrap font-sans text-sm bg-muted/50 p-4 rounded-lg">
                    {generations.prd.content.content}
                  </pre>
                </div>
              ) : (
                <EmptyState type="prd" />
              )}
            </TabsContent>

            <TabsContent value="stories" className="mt-0">
              {generations.stories ? (
                <div className="space-y-4">
                  {(generations.stories.content.stories as UserStory[]).map(
                    (story, index) => (
                      <Card key={story.id || index}>
                        <CardHeader className="pb-2">
                          <CardTitle className="text-sm font-medium">
                            {story.id || `US-${index + 1}`}: {story.title}
                          </CardTitle>
                        </CardHeader>
                        <CardContent className="space-y-2">
                          <p className="text-sm text-muted-foreground">
                            {story.description}
                          </p>
                          <Separator />
                          <div>
                            <p className="text-sm font-medium mb-1">
                              Acceptance Criteria:
                            </p>
                            <ul className="text-sm text-muted-foreground list-disc list-inside">
                              {story.acceptance_criteria?.map((criteria, i) => (
                                <li key={i}>{criteria}</li>
                              ))}
                            </ul>
                          </div>
                        </CardContent>
                      </Card>
                    ),
                  )}
                </div>
              ) : (
                <EmptyState type="stories" />
              )}
            </TabsContent>

            <TabsContent value="requirements" className="mt-0">
              {generations.requirements ? (
                <div className="prose prose-sm max-w-none dark:prose-invert">
                  <pre className="whitespace-pre-wrap font-sans text-sm bg-muted/50 p-4 rounded-lg">
                    {generations.requirements.content.content}
                  </pre>
                </div>
              ) : (
                <EmptyState type="requirements" />
              )}
            </TabsContent>

            <TabsContent value="wireframe" className="mt-0">
              {generations.wireframe ? (
                <div className="space-y-4">
                  <Card>
                    <CardHeader className="pb-2">
                      <CardTitle className="text-sm">
                        Wireframe Description
                      </CardTitle>
                    </CardHeader>
                    <CardContent>
                      <p className="text-sm text-muted-foreground">
                        {generations.wireframe.content.description}
                      </p>
                    </CardContent>
                  </Card>
                  <Card>
                    <CardHeader className="pb-2">
                      <CardTitle className="text-sm">Mermaid Diagram</CardTitle>
                    </CardHeader>
                    <CardContent>
                      <pre className="text-xs bg-muted/50 p-4 rounded-lg overflow-x-auto">
                        {generations.wireframe.content.mermaid_code}
                      </pre>
                    </CardContent>
                  </Card>
                  <Card>
                    <CardHeader className="pb-2">
                      <CardTitle className="text-sm">Preview</CardTitle>
                      <CardDescription>
                        Copy the Mermaid code to{" "}
                        <a
                          href="https://mermaid.live"
                          target="_blank"
                          rel="noopener noreferrer"
                          className="text-primary hover:underline"
                        >
                          mermaid.live
                        </a>{" "}
                        to visualize
                      </CardDescription>
                    </CardHeader>
                    <CardContent>
                      <div className="bg-muted/50 p-4 rounded-lg min-h-[200px] flex items-center justify-center">
                        <Layout className="h-8 w-8 text-muted-foreground" />
                      </div>
                    </CardContent>
                  </Card>
                </div>
              ) : (
                <EmptyState type="wireframe" />
              )}
            </TabsContent>
          </ScrollArea>
        </Tabs>
      </CardContent>
    </Card>
  );
}

function EmptyState({ type }: { type: string }) {
  const labels: Record<string, string> = {
    prd: "Product Requirements Document",
    stories: "User Stories",
    requirements: "Functional Requirements",
    wireframe: "UI Wireframe",
  };

  return (
    <div className="flex flex-col items-center justify-center py-12 text-center">
      <FileText className="h-8 w-8 text-muted-foreground mb-2" />
      <p className="text-muted-foreground text-sm">
        Click &quot;Generate&quot; to create {labels[type]}
      </p>
    </div>
  );
}
