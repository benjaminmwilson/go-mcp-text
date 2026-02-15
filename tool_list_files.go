package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var textExtensions = map[string]bool{
	".txt": true, ".md": true, ".json": true, ".yaml": true, ".yml": true,
	".toml": true, ".ini": true, ".cfg": true, ".conf": true, ".log": true,
	".csv": true, ".tsv": true, ".xml": true, ".html": true, ".htm": true,
	".css": true, ".js": true, ".ts": true, ".go": true, ".py": true,
	".rs": true, ".c": true, ".h": true, ".cpp": true, ".java": true,
	".sh": true, ".bash": true, ".zsh": true, ".fish": true,
	".sql": true, ".graphql": true, ".proto": true, ".env": true,
}

func isTextFile(name string) bool {
	return textExtensions[strings.ToLower(filepath.Ext(name))]
}

func listTextFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, e := range entries {
		if !e.IsDir() && isTextFile(e.Name()) {
			files = append(files, e.Name())
		}
	}
	return files, nil
}

func registerListFilesTool(s *server.MCPServer, sourceDir string) {
	s.AddTool(
		mcp.NewTool("list_files",
			mcp.WithDescription("List all text files available in the source directory."),
		),
		func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			files, err := listTextFiles(sourceDir)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to list files: %v", err)), nil
			}
			if len(files) == 0 {
				return mcp.NewToolResultText("No text files found in " + sourceDir), nil
			}
			return mcp.NewToolResultText(strings.Join(files, "\n")), nil
		},
	)
}
