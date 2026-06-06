// docs generates docs/openapi.json and docs/API_CONTRACTS.md from the
// single-source-of-truth in internal/apidocs.
//
// Run from the backend directory:
//
//	go run ./cmd/docs
//
// Or via go generate:
//
//	go generate ./...
package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"prd-generator/internal/apidocs"
)

func main() {
	_, thisFile, _, _ := runtime.Caller(0)
	backendRoot := filepath.Join(filepath.Dir(thisFile), "..", "..")
	docsDir := filepath.Join(backendRoot, "docs")

	if err := os.MkdirAll(docsDir, 0o755); err != nil {
		log.Fatalf("Failed to create docs dir: %v", err)
	}

	openAPIPath := filepath.Join(docsDir, "openapi.json")
	jsonBytes, err := apidocs.OpenAPIJSON()
	if err != nil {
		log.Fatalf("Failed to generate OpenAPI JSON: %v", err)
	}
	if err := os.WriteFile(openAPIPath, jsonBytes, 0o644); err != nil {
		log.Fatalf("Failed to write openapi.json: %v", err)
	}
	log.Printf("Written %s", openAPIPath)

	mdPath := filepath.Join(docsDir, "API_CONTRACTS.md")
	if err := os.WriteFile(mdPath, []byte(apidocs.Markdown()), 0o644); err != nil {
		log.Fatalf("Failed to write API_CONTRACTS.md: %v", err)
	}
	log.Printf("Written %s", mdPath)
}
