package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
)

func getBranches() ([]string, error) {
	cmd := exec.Command("git", "branch", "--format=%(refname:short)")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return strings.Split(strings.TrimSpace(string(output)), "\n"), nil
}

func main() {
	branches, err := getBranches()
	if err != nil {
		fmt.Println("Error getting branches:", err)
		return
	}

	prompt := promptui.Select{
		Label: "Select a Git Branch",
		Items: branches,
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return
	}

	cmd := exec.Command("git", "checkout", result)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error checking out branch:", err)
	}
}
