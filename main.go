package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/Codexiaoyi/ai-git/pkg/ai"
	"github.com/Codexiaoyi/ai-git/pkg/git"
	"github.com/spf13/cobra"
)

var isManual bool
var workDir string

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "AI-assisted commit, auto-generate commit message.",
	Long:  "AI-assisted commit, auto-generate commit message.",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := ai.LoadConfig()
		if err != nil {
			return fmt.Errorf("Error loading config: %v", err)
		}
		handleCommit(*config, workDir, isManual)
		return nil
	},
}

func main() {
	var rootCmd = &cobra.Command{
		Use: "ai-git",
	}
	rootCmd.AddCommand(commitCmd)
	// flags: -m, and --dir
	commitCmd.Flags().BoolVarP(&isManual, "manual", "m", false, "Commit message to use manually")
	commitCmd.Flags().StringVar(&workDir, "dir", "", "Git repository directory")
	rootCmd.Execute()
}

func handleCommit(config ai.Config, workDir string, byManual bool) {
	// Get detailed git changes information
	changes, err := git.GetChanges(workDir)
	if err != nil {
		log.Fatalf("Error getting git changes: %v", err)
	}

	// No changes to commit
	if len(changes.Modified) == 0 && len(changes.Added) == 0 && len(changes.Deleted) == 0 && len(changes.Unknown) == 0 {
		fmt.Println("No changes to commit")
		return
	}

	// Format changes for the prompt
	formattedChanges := git.FormatChangesForPrompt(changes)

	// Create prompt
	prompt := fmt.Sprintf("Generate a concise git commit message based on these changes:\n\n%s, just give me the totally summary and shortly result, you can add emojis.", formattedChanges)

	// Generate commit message using AI
	message, err := ai.GenerateCommitMessage(prompt, config)
	if err != nil {
		log.Fatalf("Error generating commit message: %v", err)
	}

	if byManual {
		// Write the AI-generated message to a temporary file for editing
		tempFile, err := os.CreateTemp("", "ai-git-commit-msg-*.txt")
		if err != nil {
			log.Fatalf("Error creating temporary file: %v", err)
		}
		defer os.Remove(tempFile.Name()) // Clean up file when done

		// Write AI-generated message to the file
		fmt.Fprintf(tempFile, "%s\n\n# AI-generated commit message. Save and close the editor to confirm the commit.\n#Or clear the file to cancel the commit.\n# Lines starting with # will be ignored.", message)
		tempFile.Close()

		// Open the temporary file in the user's default editor
		editor := os.Getenv("ai-git.editor")
		if editor == "" {
			editor = "vim" // Default to vim if EDITOR is not set
		}

		// Run the editor
		editCmd := exec.Command(editor, tempFile.Name())
		editCmd.Stdin = os.Stdin
		editCmd.Stdout = os.Stdout
		editCmd.Stderr = os.Stderr
		if err := editCmd.Run(); err != nil {
			log.Fatalf("Error opening editor: %v", err)
		}

		// Read the edited message
		editedMessageBytes, err := os.ReadFile(tempFile.Name())
		if err != nil {
			log.Fatalf("Error reading edited message: %v", err)
		}

		// Process the edited message - remove comment lines
		lines := strings.Split(string(editedMessageBytes), "\n")
		var finalLines []string
		for _, line := range lines {
			if !strings.HasPrefix(strings.TrimSpace(line), "#") {
				finalLines = append(finalLines, line)
			}
		}
		message = strings.TrimSpace(strings.Join(finalLines, "\n"))
	}

	// If the message is empty, cancel the commit
	if message == "" {
		fmt.Println("Commit message is empty. Commit cancelled.")
		return
	}

	// Add all changes
	addCmd := exec.Command("git", "add", "-A")
	addCmd.Dir = changes.WorkDir
	if err := addCmd.Run(); err != nil {
		log.Fatalf("Error executing git add: %v", err)
	}

	// Execute git commit with the edited message
	commitCmd := exec.Command("git", "commit", "-m", message)
	commitCmd.Dir = changes.WorkDir
	commitCmd.Stdout = os.Stdout
	commitCmd.Stderr = os.Stderr

	if err := commitCmd.Run(); err != nil {
		log.Fatalf("Error executing git commit: %v", err)
	}
}
