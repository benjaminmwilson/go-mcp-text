package main

import (
	"fmt"
)

func help() {
	fmt.Printf("%s v%s\n", mcpServerName, version)
	fmt.Println()
	fmt.Println("An MCP server that exposes text files as tools and resources.")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  go-mcp-text [OPTIONS]")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println("  -source <dir>   Directory to serve text files from (default: \".\")")
	fmt.Println("  -port <port>    Port to listen on in HTTP/SSE mode (default: \"8080\")")
	fmt.Println("  -stdio          Run as a stdio MCP server instead of HTTP/SSE")
	fmt.Println("  -help           Show this help message")
	fmt.Println()
	fmt.Println("TOOLS EXPOSED:")
	fmt.Println("  list_files      List all text files in the source directory")
	fmt.Println("  read_file       Read the contents of a specific file")
	fmt.Println("  search_files    Search for a string across all text files")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  go-mcp-text -source /path/to/docs")
	fmt.Println("  go-mcp-text -source /path/to/docs -port 9090")
	fmt.Println("  go-mcp-text -source /path/to/docs -stdio")
}
