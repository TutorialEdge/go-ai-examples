# MCP Server in Go

Companion code for the [Building an MCP Server in Go](https://tutorialedge.net/ai/building-an-mcp-server-in-go/)
tutorial. A minimal [Model Context Protocol](https://modelcontextprotocol.io)
server built with [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) that
exposes a single `search_tutorials` tool over stdio.

## Run it

```bash
go build -o te-mcp .
./te-mcp        # speaks MCP over stdio
```

## Quick smoke test

```bash
printf '%s\n' \
'{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"smoke","version":"1.0"}}}' \
'{"jsonrpc":"2.0","method":"notifications/initialized"}' \
'{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"search_tutorials","arguments":{"query":"agents"}}}' \
| ./te-mcp
```
