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
var newBranch bool
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

// checkoutCmd represents the checkout -b command
var checkoutCmd = &cobra.Command{
	Use:   "checkout",
	Short: "AI-assisted branch creation, auto-generate branch name.",
	Long:  "AI-assisted branch creation, auto-generate branch name based on current changes.",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := ai.LoadConfig()
		if err != nil {
			return fmt.Errorf("Error loading config: %v", err)
		}
		handleCheckout(*config, workDir, newBranch)
		return nil
	},
}

func main() {
	var rootCmd = &cobra.Command{
		Use: "ai-git",
	}
	rootCmd.AddCommand(commitCmd)
	rootCmd.AddCommand(checkoutCmd)
	// flags: -m, and --dir
	commitCmd.Flags().BoolVarP(&isManual, "manual", "m", false, "Commit message to use manually")
	commitCmd.Flags().StringVar(&workDir, "dir", "", "Git repository directory")
	checkoutCmd.Flags().BoolVarP(&newBranch, "newBranch", "b", false, "New branch name to use manually")
	checkoutCmd.Flags().StringVar(&workDir, "dir", "", "Git repository directory")
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
	prompt := fmt.Sprintf("Generate a concise git commit message based on these changes:\n\n%s, just give me the shortly commit message, you can add emojis.", formattedChanges)

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
	}

	// If the message is empty, cancel the commit
	if message == "" {
		fmt.Println("Commit message is empty. Commit cancelled.")
		return
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

func handleCheckout(config ai.Config, workDir string, newName bool) {
	if !newBranch {
		return
	}
	// Get detailed git changes information
	changes, err := git.GetChanges(workDir)
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
	checkoutCmd.Dir = changes.WorkDir
	checkoutCmd.Stdout = os.Stdout
	checkoutCmd.Stderr = os.Stderr

	if err := checkoutCmd.Run(); err != nil {
		log.Fatalf("Error executing git checkout -b: %v", err)
	}
}
