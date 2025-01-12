package main

import (
	"fmt"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
)

func formatRelativeTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	days := int(diff.Hours() / 24)
	if days == 0 {
		hours := int(diff.Hours())
		if hours == 0 {
			minutes := int(diff.Minutes())
			return fmt.Sprintf("%d minutes ago", minutes)
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else if days < 7 {
		return fmt.Sprintf("%d days ago", days)
	} else if days < 30 {
		weeks := days / 7
		return fmt.Sprintf("%d weeks ago", weeks)
	}
	months := days / 30
	return fmt.Sprintf("%d months ago", months)
}

type branchInfo struct {
	name       string
	lastCommit time.Time
}

func getBranches() ([]branchInfo, error) {
	cmd := exec.Command("git", "for-each-ref", "--sort=-committerdate", "refs/heads/", "--format=%(refname:short) %(committerdate:iso8601)")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	branches := make([]branchInfo, 0, len(lines))

	for _, line := range lines {
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			fmt.Printf("Skipping malformed line: %q\n", line)
			continue
		}
		
		commitTime, err := time.Parse("2006-01-02 15:04:05 -0700", parts[1])
		if err != nil {
			continue
		}

		branches = append(branches, branchInfo{
			name:       parts[0],
			lastCommit: commitTime,
		})
	}

	sort.Slice(branches, func(i, j int) bool {
		return branches[i].lastCommit.After(branches[j].lastCommit)
	})

	return branches, nil
}

func main() {
	branchInfos, err := getBranches()
	if err != nil {
		fmt.Println("Error getting branches:", err)
		return
	}

	// Find longest branch name for alignment
	maxLen := 0
	for _, info := range branchInfos {
		if len(info.name) > maxLen {
			maxLen = len(info.name)
		}
	}

	branchNames := make([]string, len(branchInfos))
	for i, info := range branchInfos {
		padding := strings.Repeat(" ", maxLen-len(info.name)+4) // 4 spaces minimum padding
		branchNames[i] = fmt.Sprintf("%s%s\033[2m(%s)\033[0m", info.name, padding, formatRelativeTime(info.lastCommit))
	}

	prompt := promptui.Select{
		Label: "Select a Git Branch",
		Items: branchNames,
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return
	}

	branchName := strings.Split(result, " ")[0]
	cmd := exec.Command("git", "checkout", branchName)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error checking out branch:", err)
	}
}
