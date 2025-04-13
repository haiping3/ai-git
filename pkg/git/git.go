package git

import (
	"os/exec"
	"strings"
)

// FileStatus represents the status of a file in git
type FileStatus string

const (
	Modified FileStatus = "modified"
	Added    FileStatus = "added"
	Deleted  FileStatus = "deleted"
)

// Changes represents git changes in the repository
type Changes struct {
	Modified []string            `json:"modified"`
	Added    []string            `json:"added"`
	Deleted  []string            `json:"deleted"`
	Unknown  []string            `json:"unknown"`
	Details  map[string][]string `json:"details"`
}

// GetDiff gets the git diff information for the current working directory
func GetDiff() (string, error) {
	cmd := exec.Command("git", "diff", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// GetStatus gets the git status information for the current working directory
func GetStatus() (string, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// GetChanges gets detailed information about changes in the git repository
func GetChanges() (*Changes, error) {
	// Get git diff
	diffCmd := exec.Command("git", "diff", "HEAD")
	diffOutput, err := diffCmd.Output()
	if err != nil {
		return nil, err
	}

	// Get git status
	statusCmd := exec.Command("git", "status", "--porcelain")
	statusOutput, err := statusCmd.Output()
	if err != nil {
		return nil, err
	}

	// Initialize changes
	changes := &Changes{
		Modified: make([]string, 0),
		Added:    make([]string, 0),
		Deleted:  make([]string, 0),
		Unknown:  make([]string, 0),
		Details:  make(map[string][]string),
	}

	// Process status output to categorize files
	statusLines := strings.Split(strings.TrimSpace(string(statusOutput)), "\n")
	for _, line := range statusLines {
		if line == "" {
			continue
		}

		status := strings.TrimSpace(line[0:2])
		file := line[3:]

		if strings.Contains(status, "M") {
			changes.Modified = append(changes.Modified, file)
		}
		if strings.Contains(status, "A") {
			changes.Added = append(changes.Added, file)
		}
		if strings.Contains(status, "D") {
			changes.Deleted = append(changes.Deleted, file)
		}
		if strings.Contains(status, "??") {
			changes.Unknown = append(changes.Unknown, file)
		}
	}

	// Process diff output to get detailed changes
	diffLines := strings.Split(string(diffOutput), "\n")
	var currentFile string

	for _, line := range diffLines {
		if strings.HasPrefix(line, "diff --git") {
			parts := strings.Split(line, " b/")
			if len(parts) > 1 {
				currentFile = parts[1]
				changes.Details[currentFile] = make([]string, 0)
			}
		} else if (strings.HasPrefix(line, "+") || strings.HasPrefix(line, "-")) &&
			!strings.HasPrefix(line, "+++") && !strings.HasPrefix(line, "---") {
			if currentFile != "" {
				changes.Details[currentFile] = append(changes.Details[currentFile], line)
			}
		}
	}

	return changes, nil
}

// FormatChangesForPrompt converts the Changes structure to a formatted string for use in AI prompts
func FormatChangesForPrompt(changes *Changes) string {
	var sb strings.Builder

	sb.WriteString("Git Changes Summary:\n\n")

	// Add summary information
	if len(changes.Modified) > 0 {
		sb.WriteString("Modified files:\n")
		for _, file := range changes.Modified {
			sb.WriteString("- " + file + "\n")
		}
		sb.WriteString("\n")
	}

	if len(changes.Added) > 0 {
		sb.WriteString("Added files:\n")
		for _, file := range changes.Added {
			sb.WriteString("- " + file + "\n")
		}
		sb.WriteString("\n")
	}

	if len(changes.Deleted) > 0 {
		sb.WriteString("Deleted files:\n")
		for _, file := range changes.Deleted {
			sb.WriteString("- " + file + "\n")
		}
		sb.WriteString("\n")
	}

	if len(changes.Unknown) > 0 {
		sb.WriteString("Unknown files:\n")
		for _, file := range changes.Unknown {
			sb.WriteString("- " + file + "\n")
		}
		sb.WriteString("\n")
	}

	// Add detailed changes
	if len(changes.Details) > 0 {
		sb.WriteString("Detailed Changes:\n\n")
		for file, lines := range changes.Details {
			sb.WriteString("File: " + file + "\n")
			for _, line := range lines {
				sb.WriteString(line + "\n")
			}
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
