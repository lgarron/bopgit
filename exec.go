package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/logrusorgru/aurora"
)

func gitExecCommand(args ...string) *exec.Cmd {
	return exec.Command("git", args...)
}

func runGitCommand(args ...string) string {
	output, err := gitExecCommand(args...).Output()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return strings.TrimSuffix(string(output), "\n")
}

func branchMustExist(branch gitBranch) {
	if !doesBranchExist(branch) {
		fmt.Errorf("Branch does not exist: ",
			aurora.Bold(branch),
		)
		showHelpAndExit()
	}
}
