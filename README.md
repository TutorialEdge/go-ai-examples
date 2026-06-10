# Go AI Examples

Companion code for the [TutorialEdge](https://tutorialedge.net) AI-in-Go tutorials.
Clone it, `cd` into any example, and run it — each one is small, self-contained,
and matches the written tutorial (and video) step for step.

## How this repo is organised

Every top-level directory is **one tutorial and its own Go module**, so each
example is independently runnable and they never fight over dependencies. Pick
the one you want and follow its steps below.

| Example | What you build | Tutorial |
| --- | --- | --- |
| [`adk-go-agent/`](./adk-go-agent) | Your first AI agent with Google's Agent Development Kit (Gemini + agent + session + runner, streaming the reply) | [Build Your First AI Agent in Go with the ADK](https://tutorialedge.net/ai/build-your-first-ai-agent-in-go-with-adk/) |
| [`mcp-server-go/`](./mcp-server-go) | A Model Context Protocol server exposing a `search_tutorials` tool over stdio | [Building an MCP Server in Go](https://tutorialedge.net/ai/building-an-mcp-server-in-go/) |
| [`mcp-server-security-go/`](./mcp-server-security-go) | Harden an MCP server: break it with path traversal, then add path confinement, input limits, output hygiene, rate limiting, and tool annotations | [MCP Server Security in Go](https://tutorialedge.net/ai/mcp-server-security-in-go/) |

### Also part of the series

- **Getting Started with the Claude API in Go** — a full 12-part course lives in
  its own repo: [TutorialEdge/claude-api-go](https://github.com/TutorialEdge/claude-api-go)
  ([tutorial](https://tutorialedge.net/ai/getting-started-with-claude-api-in-go/)).

## Prerequisites

- **Go 1.23 or newer** — some examples use range-over-function iterators (Go 1.23+).
- An **API key** for whichever provider an example uses (see each example below).
  Keep keys in your environment or a secrets manager — never commit them.

## adk-go-agent

Your first agent with the [ADK for Go](https://github.com/google/adk-go). Needs a
Google AI Studio API key ([get one free](https://aistudio.google.com/apikey)).

```bash
cd adk-go-agent
export GOOGLE_API_KEY="your-key-here"
go run .
```

It sends a question to a Gemini-backed agent and streams the answer to your
terminal. Change the `Instruction` field in `main.go` to change the agent's whole
personality — that's the heart of the tutorial.

## mcp-server-go

A minimal MCP server built with [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go),
speaking the Model Context Protocol over stdio. No API key required.

```bash
cd mcp-server-go
go build -o te-mcp .
./te-mcp        # speaks MCP over stdio
```

See [`mcp-server-go/README.md`](./mcp-server-go/README.md) for a one-liner smoke
test and how to wire it into a client like Claude.

## mcp-server-security-go

The security follow-up to `mcp-server-go`: a deliberately vulnerable `read_doc`
tool that the tutorial attacks with path traversal, then hardens step by step.
No API key required.

```bash
cd mcp-server-security-go
go build -o te-mcp-secure .
./te-mcp-secure   # speaks MCP over stdio
```

See [`mcp-server-security-go/README.md`](./mcp-server-security-go/README.md) for
the attack payload and the three hardening steps.

## License

[MIT](./LICENSE) © Elliot Forbes / TutorialEdge.
