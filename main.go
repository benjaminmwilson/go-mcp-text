package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
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

func main() {
	var sourceDir string
	var port string

	flag.StringVar(&sourceDir, "source", ".", "Directory to serve text files from")
	flag.StringVar(&port, "port", "8080", "Port to listen on")
	flag.Parse()

	// Resolve to absolute path and validate.
	absDir, err := filepath.Abs(sourceDir)
	if err != nil {
		log.Fatalf("failed to resolve source directory: %v", err)
	}
	info, err := os.Stat(absDir)
	if err != nil || !info.IsDir() {
		log.Fatalf("source directory does not exist or is not a directory: %s", absDir)
	}
	sourceDir = absDir

	s := server.NewMCPServer(
		"file-reader",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, false),
	)

	// Tool: list_files — enumerate text files in the source directory.
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

	// Tool: read_file — read a named text file from the source directory.
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

	// Tool: search_files — search all text files for a keyword.
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

	addr := ":" + port
	fmt.Printf("MCP file-reader server starting\n")
	fmt.Printf("  Source directory : %s\n", sourceDir)
	fmt.Printf("  Listening on     : http://localhost%s/mcp\n", addr)

	httpServer := server.NewStreamableHTTPServer(s, server.WithStateLess(true))
	if err := httpServer.Start(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
