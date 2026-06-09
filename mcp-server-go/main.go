package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// tutorial is a single piece of TutorialEdge content the server can search.
type tutorial struct {
	Title string
	URL   string
	Tags  []string
}

// catalog is a tiny in-memory stand-in for a real content index. In a
// production server you'd query a database or a search API here instead.
var catalog = []tutorial{
	{
		Title: "Getting Started with the Claude API in Go",
		URL:   "https://tutorialedge.net/ai/getting-started-with-claude-api-in-go/",
		Tags:  []string{"ai", "go", "anthropic"},
	},
	{
		Title: "Building AI Agents in Go",
		URL:   "https://tutorialedge.net/ai/building-ai-agents-in-go/",
		Tags:  []string{"ai", "go", "agents"},
	},
	{
		Title: "Building RAG Applications in Go",
		URL:   "https://tutorialedge.net/ai/building-rag-applications-in-go/",
		Tags:  []string{"ai", "go", "rag"},
	},
	{
		Title: "Calling Ollama from Go",
		URL:   "https://tutorialedge.net/ai/calling-ollama-from-go/",
		Tags:  []string{"ai", "go", "ollama"},
	},
}

func main() {
	// Create the server. The name and version are reported to clients
	// during the MCP handshake.
	s := server.NewMCPServer(
		"tutorialedge",
		"1.0.0",
	)

	// Describe a single tool. The description and parameter schema are how
	// the LLM decides when and how to call it, so be specific.
	tool := mcp.NewTool("search_tutorials",
		mcp.WithDescription("Search TutorialEdge content for tutorials matching a query"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description(`The search term, e.g. "go agents"`),
		),
	)

	// Register the tool with its handler, then serve over stdio.
	s.AddTool(tool, handleSearch)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("server error: %v\n", err)
	}
}

// handleSearch is invoked whenever a client calls the search_tutorials tool.
func handleSearch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, err := request.RequireString("query")
	if err != nil {
		// A bad request is a tool error, not a transport error — report it
		// back to the model so it can correct the call.
		return mcp.NewToolResultError(err.Error()), nil
	}

	matches := search(query)
	if len(matches) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("No tutorials found for %q.", query)), nil
	}

	var b strings.Builder
	fmt.Fprintf(&b, "Found %d tutorial(s) for %q:\n\n", len(matches), query)
	for _, t := range matches {
		fmt.Fprintf(&b, "- %s\n  %s\n", t.Title, t.URL)
	}

	return mcp.NewToolResultText(b.String()), nil
}

// search does a naive case-insensitive match across titles and tags.
func search(query string) []tutorial {
	q := strings.ToLower(query)
	var out []tutorial
	for _, t := range catalog {
		if strings.Contains(strings.ToLower(t.Title), q) {
			out = append(out, t)
			continue
		}
		for _, tag := range t.Tags {
			if strings.Contains(strings.ToLower(tag), q) {
				out = append(out, t)
				break
			}
		}
	}
	return out
}
