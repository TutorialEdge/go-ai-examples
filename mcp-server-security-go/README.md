# MCP Server Security in Go

Companion code for the [MCP Server Security in Go — Hardening Your Server](https://tutorialedge.net/ai/mcp-server-security-in-go/)
tutorial. A minimal [Model Context Protocol](https://modelcontextprotocol.io)
server built with [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) that
exposes a `read_doc` tool over stdio — and starts **deliberately vulnerable**
to path traversal so the tutorial can break it, then harden it.

The hardening steps are marked `STEP 1`–`STEP 5` in `main.go`:

1. `safeJoin` — confine paths to the docs directory (kills traversal)
2. Input limits — length cap + `.md` extension allowlist
3. Output hygiene — response size cap, and label served content as untrusted data
4. Bound time and volume — context timeout + rate limiter, so a surge
   (malicious or otherwise) gets a tool error instead of service
5. Tool annotations — declare `read_doc` read-only/idempotent; flag
   destructive tools so clients put a human in the loop

## Run it

```bash
go build -o te-mcp-secure .
./te-mcp-secure   # speaks MCP over stdio
```

## Reproduce the attack

```bash
printf '%s\n' \
'{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"smoke","version":"1.0"}}}' \
'{"jsonrpc":"2.0","method":"notifications/initialized"}' \
'{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"read_doc","arguments":{"path":"../../../../etc/passwd"}}}' \
| ./te-mcp-secure
```

Before the hardening steps this leaks `/etc/passwd`; after, it returns a tool
error and the process stays up.
