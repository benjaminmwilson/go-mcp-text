package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	var sourceDir string
	var port string
	var useStdio bool

	flag.StringVar(&sourceDir, "source", ".", "Directory to serve text files from")
	flag.StringVar(&port, "port", "8080", "Port to listen on")
	flag.BoolVar(&useStdio, "stdio", false, "Run as a stdio MCP server instead of HTTP/SSE")
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

	registerListFilesTool(s, sourceDir)
	registerReadFileTool(s, sourceDir)
	registerSearchFilesTool(s, sourceDir)

	registerResourceTemplates(s, sourceDir)

	if useStdio {
		if err := server.ServeStdio(s); err != nil {
			log.Fatalf("stdio server error: %v", err)
		}
	} else {
		addr := ":" + port
		fmt.Printf("MCP file-reader server starting\n")
		fmt.Printf("  Source directory : %s\n", sourceDir)
		fmt.Printf("  Listening on     : http://localhost%s/mcp\n", addr)

		httpServer := server.NewStreamableHTTPServer(s, server.WithStateLess(true))
		if err := httpServer.Start(addr); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}
}
