package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// searchFiles searches all text files in dir for lines containing query (case-insensitive).
// Returns a slice of "filename:linenum: line" strings.
func searchFiles(dir, query string) ([]string, error) {
	files, err := listTextFiles(dir)
	if err != nil {
		return nil, err
	}
	lower := strings.ToLower(query)
	var matches []string
	for _, name := range files {
		abs := filepath.Join(dir, name)
		f, err := os.Open(abs)
		if err != nil {
			continue
		}
		scanner := bufio.NewScanner(f)
		lineNum := 0
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()
			if strings.Contains(strings.ToLower(line), lower) {
				matches = append(matches, fmt.Sprintf("%s:%d: %s", name, lineNum, line))
			}
		}
		f.Close()
	}
	return matches, nil
}

func registerSearchFilesTool(s *server.MCPServer, sourceDir string) {
	s.AddTool(
		mcp.NewTool("search_files",
			mcp.WithDescription("Search all text files in the source directory for lines matching a keyword (case-insensitive)."),
			mcp.WithString("query",
				mcp.Required(),
				mcp.Description("Keyword or phrase to search for."),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			query, err := req.RequireString("query")
			if err != nil {
				return mcp.NewToolResultError("query parameter is required"), nil
			}
			matches, err := searchFiles(sourceDir, query)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to search files: %v", err)), nil
			}
			if len(matches) == 0 {
				return mcp.NewToolResultText("No matches found."), nil
			}
			return mcp.NewToolResultText(strings.Join(matches, "\n")), nil
		},
	)
}
