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

// safeRead reads a file, ensuring it stays within dir to prevent path traversal.
func safeRead(dir, filename string) (string, error) {
	clean := filepath.Clean(filename)
	if strings.HasPrefix(clean, "..") {
		return "", fmt.Errorf("invalid filename: path traversal not allowed")
	}
	abs, err := filepath.Abs(filepath.Join(dir, clean))
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(abs, dir+string(filepath.Separator)) && abs != dir {
		return "", fmt.Errorf("access denied: file is outside source directory")
	}
	data, err := os.ReadFile(abs)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func registerReadFileTool(s *server.MCPServer, sourceDir string) {
	s.AddTool(
		mcp.NewTool("read_file",
			mcp.WithDescription("Read the contents of a text file from the source directory."),
			mcp.WithString("filename",
				mcp.Required(),
				mcp.Description("Name of the file to read (relative to source directory, no path traversal)."),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			filename, err := req.RequireString("filename")
			if err != nil {
				return mcp.NewToolResultError("filename parameter is required"), nil
			}
			content, err := safeRead(sourceDir, filename)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to read file: %v", err)), nil
			}
			return mcp.NewToolResultText(content), nil
		},
	)
}
