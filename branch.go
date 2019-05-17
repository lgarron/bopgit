package main

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
)

type gitBranch = string

func hash(ref string) string {
	return runGitCommand("show-ref", "-d", "-s", ref)
}

func currentBranch() string {
	return runGitCommand("rev-parse", "--abbrev-ref", "HEAD")
}

func doesBranchExist(branch gitBranch) bool {
	return isGitCommandExitCodeZero("rev-parse", "--verify", branch)
}

func branchMustExist(branch gitBranch) {
	if !doesBranchExist(branch) {
		fmt.Printf("Branch does not exist: %s",
			aurora.Bold(branch),
		)
		showHelp()
		os.Exit(1)
	}
}

func doesBranchContain(branch gitBranch, ref string) bool {
	return isGitCommandExitCodeZero("merge-base", "--is-ancestor", ref, branch)
}

func branchMustContain(branch gitBranch, ref string) {
	if !doesBranchContain(branch, ref) {
		fmt.Errorf("Branch %s does not contain expected ref: %s",
			aurora.Bold(branch),
			aurora.Bold(ref),
		)
		showHelp()
		os.Exit(1)
	}
}
