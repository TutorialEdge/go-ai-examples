package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/genai"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
)

func main() {
	ctx := context.Background()

	// The model — the brain the agent reasons with.
	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatal(err)
	}

	// The agent — a name, what it's for, and the instruction that shapes it.
	assistant, err := llmagent.New(llmagent.Config{
		Name:        "assistant",
		Model:       model,
		Description: "A helpful Go programming assistant.",
		Instruction: "You are a concise, friendly assistant for Go developers.",
	})
	if err != nil {
		log.Fatal(err)
	}

	// A session holds the conversation; the runner drives the agent.
	sessions := session.InMemoryService()
	created, err := sessions.Create(ctx, &session.CreateRequest{
		AppName: "app",
		UserID:  "user",
	})
	if err != nil {
		log.Fatal(err)
	}

	r, err := runner.New(runner.Config{
		AppName:        "app",
		Agent:          assistant,
		SessionService: sessions,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Send a message and stream the reply.
	msg := genai.NewContentFromText("What's a goroutine?", genai.RoleUser)
	for event, err := range r.Run(ctx, "user", created.Session.ID(), msg, agent.RunConfig{}) {
		if err != nil {
			log.Fatal(err)
		}
		if event.Content != nil {
			for _, part := range event.Content.Parts {
				fmt.Print(part.Text)
			}
		}
	}
	fmt.Println()
}
