package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerResourceTemplates(s *server.MCPServer, sourceDir string) {
	// Resource template: file://{filename} — expose files as MCP resources.
	s.AddResourceTemplate(
		mcp.NewResourceTemplate(
			"file://{filename}",
			"Text File",
			mcp.WithTemplateDescription("Read a text file from the source directory via its URI."),
			mcp.WithTemplateMIMEType("text/plain"),
		),
		func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			filename := strings.TrimPrefix(req.Params.URI, "file://")
			content, err := safeRead(sourceDir, filename)
			if err != nil {
				return nil, fmt.Errorf("failed to read %s: %w", filename, err)
			}
			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      req.Params.URI,
					MIMEType: "text/plain",
					Text:     content,
				},
			}, nil
		},
	)
}
