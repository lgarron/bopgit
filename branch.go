package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/logrusorgru/aurora"
)

type Branch struct {
	name string
}

func (b Branch) String() string {
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

func newBranch(branchName string) Branch {
	branchNameMustExist(branchName)
	return Branch{
		name: branchName,
	}
}

func hash(ref string) string {
	return runGitCommand("show-ref", "--heads", "-s", ref)
}

func currentBranch() Branch {
	return newBranch(runGitCommand("rev-parse", "--abbrev-ref", "HEAD"))
}

func doesBranchContain(branch Branch, ref string) bool {
	return isGitCommandExitCodeZero("merge-base", "--is-ancestor", ref, branch.name)
}

func branchMustContain(branch Branch, ref string) {
	if !doesBranchContain(branch, ref) {
		fmt.Errorf("Branch %s does not contain expected ref: %s",
			aurora.Bold(branch),
			aurora.Bold(ref),
		)
		showHelp()
		os.Exit(1)
	}
}

func rebaseOnto(newbase string, upstream string, root Branch) {
	runGitCommand("rebase", "--onto", newbase, upstream, root.name)
}

func numCommitsAhead(branch string, comparison string) int {
	s := runGitCommand("rev-list", "--left-only", "--count", fmt.Sprintf(
		"%s...%s",
		branch,
		comparison,
	))
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	return i
}
