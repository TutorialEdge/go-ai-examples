// Starter code for "MCP Server Security in Go — Hardening Your Server".
//
// This server exposes a single read_doc tool that serves markdown files from
// ./docs. It starts deliberately VULNERABLE — on camera we attack it with
// path traversal, then harden it step by step:
//
//	STEP 1: implement safeJoin to confine paths to docsRoot
//	STEP 2: add input limits (length cap, .md extension allowlist)
//	STEP 3: cap response size and label output as untrusted data
//	STEP 4: bound time and volume (context timeout + rate limiter)
//	STEP 5: annotate the tool (read-only, idempotent) so clients can
//	        gate destructive tools behind a human
package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// docsRoot is the only directory this server should ever read from.
const docsRoot = "docs"

// maxDocBytes caps how much file content one tool call may return.
const maxDocBytes = 64 * 1024

func main() {
	s := server.NewMCPServer(
		"tutorialedge-secure",
		"1.0.0",
	)

	// STEP 5 (on camera): annotate the tool — read_doc is read-only and
	// idempotent, so say so. Destructive tools get
	// mcp.WithDestructiveHintAnnotation(true) so clients ask a human first.
	tool := mcp.NewTool("read_doc",
		mcp.WithDescription("Read a markdown document from the docs directory"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description(`Relative path of the doc to read, e.g. "welcome.md"`),
		),
	)

	s.AddTool(tool, handleReadDoc)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("server error: %v\n", err)
	}
}

// handleReadDoc is invoked whenever a client calls the read_doc tool.
func handleReadDoc(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path, err := request.RequireString("path")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// STEP 4 (on camera): bound time and volume —
	//   - ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	//   - a rate.Limiter (30 calls/min): exhausted → tool error, not service

	// STEP 2 (on camera): input limits go here —
	//   - reject paths longer than 255 chars or empty
	//   - allowlist: only ".md" files may be served

	// VULNERABLE: trusts the model's path argument completely.
	// "../../../../etc/passwd" walks straight out of docsRoot.
	// STEP 1 (on camera): replace this with safeJoin(docsRoot, path).
	full := filepath.Join(docsRoot, path)

	data, err := os.ReadFile(full)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("could not read %q: %v", path, err)), nil
	}

	// STEP 3 (on camera): cap len(data) at maxDocBytes and wrap the content
	// in a fence labelling it as untrusted file data, not instructions.
	return mcp.NewToolResultText(string(data)), nil
}

// safeJoin confines path to root, rejecting absolute paths and traversal.
// STEP 1 (on camera): implement this — filepath.Clean("/"+path) collapses
// any "..", then an Abs + prefix check guarantees the result stays inside
// root. See SCRIPT.md for the full implementation.
func safeJoin(root, path string) (string, error) {
	return "", fmt.Errorf("not implemented yet — written on camera")
}
