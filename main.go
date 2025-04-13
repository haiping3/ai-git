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

var (
	config *ai.Config
)

func main() {
	cfg, err := ai.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}
	config = cfg
	var rootCmd = &cobra.Command{
		Use:                "ai-git [command]",
		Short:              "AI-assisted git commands",
		Long:               "AI-git is a git wrapper with AI capabilities for certain commands",
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(args)
			if len(args) == 0 {
				cmd.Help()
				return
			}
			var newBranch, message string
			cmd.Flags().StringVarP(&message, "message", "m", "", "Auto generate commit message")
			cmd.Flags().StringVarP(&newBranch, "newBranch", "b", "", "Auto generate branch name")
			if err := cmd.Flags().Parse(args); err == nil {
				// Handle specific commands
				switch args[0] {
				case "commit":
					fmt.Println(cmd.Flags().Changed("message"))
					if message == "" {
						handleCommit(*config)
						return
					}
				case "checkout":
					if newBranch == "" {
						handleCheckout(*config)
						return
					}
				}
			}
			// Fallback to standard git
			gitCmd := exec.Command("git", args...)
			gitCmd.Stdin = os.Stdin
			gitCmd.Stdout = os.Stdout
			gitCmd.Stderr = os.Stderr
			if err := gitCmd.Run(); err != nil {
				if exitError, ok := err.(*exec.ExitError); ok {
					os.Exit(exitError.ExitCode())
				}
				fmt.Fprintf(os.Stderr, "Error executing git command: %v\n", err)
				os.Exit(1)
			}
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func handleCommit(config ai.Config) {
	// Get detailed git changes information
	changes, err := git.GetChanges()
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
	prompt := fmt.Sprintf("Generate a concise git commit message based on these changes:\n\n%s, just give me the shortly commit message, you can add emojis.", formattedChanges)

	// Generate commit message using AI
	message, err := ai.GenerateCommitMessage(prompt, config)
	if err != nil {
		log.Fatalf("Error generating commit message: %v", err)
	}

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
	editor := os.Getenv("AI_GIT_EDITOR")
	if editor == "" {
		editor = os.Getenv("EDITOR")
		if editor == "" {
			editor = "vim" // Default to vim if no editor is set
		}
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

	// If the message is empty, cancel the commit
	if message == "" {
		fmt.Println("Commit message is empty. Commit cancelled.")
		return
	}

	// Execute git commit with the edited message
	commitCmd := exec.Command("git", "commit", "-m", message)
	commitCmd.Stdout = os.Stdout
	commitCmd.Stderr = os.Stderr

	if err := commitCmd.Run(); err != nil {
		log.Fatalf("Error executing git commit: %v", err)
	}
}

func handleCheckout(config ai.Config) {
	// Get detailed git changes information
	changes, err := git.GetChanges()
	if err != nil {
		log.Fatalf("Error getting git changes: %v", err)
	}

	// Format changes for the prompt
	formattedChanges := git.FormatChangesForPrompt(changes)

	// Create prompt
	prompt := fmt.Sprintf("Generate a concise git branch name based on these changes:\n\n%s\n\nPlease generate a branch name that follows git branch naming conventions (lowercase, hyphen-separated, descriptive). Just give me the branch name, no explanation needed.", formattedChanges)

	// Generate branch name using AI
	branchName, err := ai.GenerateBranchName(prompt, config)
	if err != nil {
		log.Fatalf("Error generating branch name: %v", err)
	}

	// Clean up the branch name
	branchName = strings.TrimSpace(branchName)
	branchName = strings.ToLower(branchName)
	branchName = strings.ReplaceAll(branchName, " ", "-")

	// Write the AI-generated branch name to a temporary file for editing
	tempFile, err := os.CreateTemp("", "ai-git-branch-name-*.txt")
	if err != nil {
		log.Fatalf("Error creating temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up file when done

	// Write AI-generated branch name to the file
	fmt.Fprintf(tempFile, "%s\n\n# AI-generated branch name. Save and close the editor to confirm.\n# Or clear the file to cancel.\n# Lines starting with # will be ignored.", branchName)
	tempFile.Close()

	// Open the temporary file in the user's default editor
	editor := os.Getenv("AI_GIT_EDITOR")
	if editor == "" {
		editor = os.Getenv("EDITOR")
		if editor == "" {
			editor = "vim" // Default to vim if no editor is set
		}
	}

	// Run the editor
	editCmd := exec.Command(editor, tempFile.Name())
	editCmd.Stdin = os.Stdin
	editCmd.Stdout = os.Stdout
	editCmd.Stderr = os.Stderr
	if err := editCmd.Run(); err != nil {
		log.Fatalf("Error opening editor: %v", err)
	}

	// Read the edited branch name
	editedNameBytes, err := os.ReadFile(tempFile.Name())
	if err != nil {
		log.Fatalf("Error reading edited branch name: %v", err)
	}

	// Process the edited branch name - remove comment lines
	lines := strings.Split(string(editedNameBytes), "\n")
	var finalLines []string
	for _, line := range lines {
		if !strings.HasPrefix(strings.TrimSpace(line), "#") {
			finalLines = append(finalLines, line)
		}
	}
	branchName = strings.TrimSpace(strings.Join(finalLines, "\n"))

	// If the branch name is empty, cancel the operation
	if branchName == "" {
		fmt.Println("Branch name is empty. Operation cancelled.")
		return
	}

	// Execute git checkout -b with the branch name
	checkoutCmd := exec.Command("git", "checkout", "-b", branchName)
	checkoutCmd.Stdout = os.Stdout
	checkoutCmd.Stderr = os.Stderr

	if err := checkoutCmd.Run(); err != nil {
		log.Fatalf("Error executing git checkout -b: %v", err)
	}
}
