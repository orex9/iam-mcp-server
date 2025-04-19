test-local:
	npx @modelcontextprotocol/inspector build/iam-mcp-server stdio
build:
	go build -o build/iam-mcp-server main.go
