package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/ayilin/ai-git/pkg/ai"
	"github.com/ayilin/ai-git/pkg/git"
)

func main() {
	// Parse command line flags
	workDir := flag.String("dir", "", "Git repository directory")
	flag.Parse()

	// Load configuration
	config, err := ai.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Check for command
	args := flag.Args()
	if len(args) > 0 && args[0] == "commit" {
		handleCommit(*config, *workDir)
	} else {
		fmt.Println("Usage: ai-git [--dir path] commit")
	}
}

func handleCommit(config ai.Config, workDir string) {
	// Get detailed git changes information
	changes, err := git.GetChanges(workDir)
	if err != nil {
		log.Fatalf("Error getting git changes: %v", err)
	}

	// No changes to commit
	if len(changes.Modified) == 0 && len(changes.Added) == 0 && len(changes.Deleted) == 0 {
		fmt.Println("No changes to commit")
		return
	}

	// Format changes for the prompt
	formattedChanges := git.FormatChangesForPrompt(changes)

	// Create prompt
	prompt := fmt.Sprintf("Generate a concise git commit message based on these changes:\n\n%s, just give me the totally summary and shortly result", formattedChanges)

	// Generate commit message using AI
	message, err := ai.GenerateCommitMessage(prompt, config)
	if err != nil {
		log.Fatalf("Error generating commit message: %v", err)
	}

	addCmd := exec.Command("git", "add", "-A")
	if err := addCmd.Run(); err != nil {
		log.Fatalf("Error executing git add: %v", err)
	}
	// Execute git commit
	commitCmd := exec.Command("git", "commit", "-m", message)
	commitCmd.Dir = changes.WorkDir
	commitCmd.Stdout = os.Stdout
	commitCmd.Stderr = os.Stderr

	if err := commitCmd.Run(); err != nil {
		log.Fatalf("Error executing git commit: %v", err)
	}
}
