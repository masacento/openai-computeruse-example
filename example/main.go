package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	cu "github.com/masacento/openai-computeruse-example"
)

func main() {
	if os.Getenv("OPENAI_API_KEY") == "" {
		log.Fatal("OPENAI_API_KEY environment variable is not set")
	}
	url := flag.String("url", "https://duckduckgo.com/", "Initial URL")
	prompt := flag.String("prompt", "Find out the winner of the Academy Award for Best Picture in 2025 and tell me the title.", "Instruction to execute")
	maxturns := flag.Int("maxturns", 16, "Maximum number of turns (optional)")
	timeout := flag.String("timeout", "3m", "Timeout duration (optional)")
	flag.Parse()

	to, err := time.ParseDuration(*timeout)
	if err != nil {
		log.Fatalf("invalid timeout: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), to)
	defer cancel()

	fmt.Println("Prompt:", *prompt)
	fmt.Println("URL   :", *url)

	err = cu.BrowserUse(ctx, *url, *prompt, *maxturns)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Println("Done")
}
