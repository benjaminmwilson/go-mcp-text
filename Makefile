build:
	go build -ldflags="-s -w" -o go-mcp-text ./...

vet:
	go vet ./...

clean:
	rm -f go-mcp-text go-mcp-text-stripped
