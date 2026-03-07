# go-mcp-text

An MCP (Model Context Protocol) server written in Go that exposes a directory of text files as tools and resources. Supports both HTTP/SSE and stdio transports.

## Tools

| Tool | Description |
|------|-------------|
| `list_files` | List all text files in the source directory |
| `read_file` | Read the contents of a specific file |
| `search_files` | Case-insensitive keyword search across all text files |

Files are also accessible as MCP resources via the `file://{filename}` URI template.


---

## Running the server

### Prerequisites

- Go 1.25+

### Build

```sh
make build
```

This produces a `go-mcp-text` binary in the project root.

### HTTP/SSE mode (default)

```sh
./go-mcp-text -source /path/to/your/docs
```

The server starts on port 8080 by default. Connect your MCP client to:

```
http://localhost:8080/mcp
```

Use `-port` to change the port:

```sh
./go-mcp-text -source /path/to/your/docs -port 9090
```

### stdio mode (Claude Desktop / CLI clients)

```sh
./go-mcp-text -source /path/to/your/docs -stdio
```

### Options

```
-source <dir>   Directory to serve text files from (default: ".")
-port <port>    Port to listen on in HTTP/SSE mode (default: "8080")
-stdio          Run as a stdio MCP server instead of HTTP/SSE
-help           Show help
```

---

## Claude Desktop integration

Add the following to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "go-mcp-text": {
      "command": "/path/to/go-mcp-text",
      "args": ["-source", "/path/to/your/docs", "-stdio"]
    }
  }
}
```

---

## Claude Code integration

Use the `claude mcp add` command to register the server. Changes are saved to `.mcp.json` in the project root (shared with your team) when using `--scope project`, or to `~/.claude.json` for personal/local use (the default).

### stdio transport

```sh
claude mcp add --transport stdio go-mcp-text -- /path/to/go-mcp-text -source /path/to/your/docs -stdio
```

### HTTP transport

First start the server in HTTP mode:

```sh
./go-mcp-text -source /path/to/your/docs -port 8080
```

Then register it:

```sh
claude mcp add --transport http go-mcp-text http://localhost:8080/mcp
```

### Resulting `.mcp.json` / `.claude.json`

**stdio:**

```json
{
  "mcpServers": {
    "go-mcp-text": {
      "type": "stdio",
      "command": "/path/to/go-mcp-text",
      "args": ["-source", "/path/to/your/docs", "-stdio"]
    }
  }
}
```

**HTTP:**

```json
{
  "mcpServers": {
    "go-mcp-text": {
      "type": "http",
      "url": "http://localhost:8080/mcp"
    }
  }
}
```


---

## Development

### Prerequisites

- Go 1.25+

### Clone and install dependencies

```sh
git clone https://github.com/your-org/go-mcp-text.git
cd go-mcp-text
go mod download
```

### Build

```sh
make build
```

### Vet

```sh
make vet
```

### Clean

```sh
make clean
```

---

## Quick Test

This command sends the MCP handshake (`initialize` and `initialized`) and then lists the tools the MCP server offers. It needs `jq` to pretty print the results:

```bash
{ echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"0"}}}'; \
echo '{"jsonrpc":"2.0","method":"notifications/initialized","params":{}}'; \
echo '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}'; } \
| ./go-mcp-text --stdio | jq .
```


---

## Dependencies

- [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) — MCP server framework for Go
