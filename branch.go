package main

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
)

type gitBranch struct {
	name string
}

func (b gitBranch) String() string {
	return fmt.Sprintf("⌥ %s", b.name)
}

func doesBranchNameExist(branchName string) bool {
	return isGitCommandExitCodeZero("rev-parse", "--verify", branchName)
}

func branchNameMustExist(branch string) {
	if !doesBranchNameExist(branch) {
		fmt.Printf("Branch does not exist: %s",
			aurora.Bold(branch),
		)
		showHelp()
		os.Exit(1)
	}
}

func newGitBranch(branchName string) gitBranch {
	branchNameMustExist(branchName)
	return gitBranch{
		name: branchName,
	}
}

func hash(ref string) string {
	return runGitCommand("show-ref", "-d", "-s", ref)
}

func currentBranch() gitBranch {
	return newGitBranch(runGitCommand("rev-parse", "--abbrev-ref", "HEAD"))
}

func doesBranchContain(branch gitBranch, ref string) bool {
	return isGitCommandExitCodeZero("merge-base", "--is-ancestor", ref, branch.name)
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

func rebaseOnto(newbase string, upstream string, root gitBranch) {
	runGitCommand("rebase", "--into", newbase, upstream, root.name)
}
