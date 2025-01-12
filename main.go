package main

import (
	"fmt"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
)

type branchInfo struct {
	name       string
	lastCommit time.Time
}

func getBranches() ([]branchInfo, error) {
	cmd := exec.Command("git", "for-each-ref", "--sort=-committerdate", "refs/heads/", "--format=%(refname:short)%09%(committerdate:iso8601)")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	branches := make([]branchInfo, 0, len(lines))

	for _, line := range lines {
		parts := strings.Split(line, "\t")
		if len(parts) != 2 {
			continue
		}
		
		commitTime, err := time.Parse(time.RFC3339, parts[1])
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

	branchNames := make([]string, len(branchInfos))
	for i, info := range branchInfos {
		branchNames[i] = fmt.Sprintf("%s (%s)", info.name, info.lastCommit.Format("2006-01-02 15:04:05"))
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
