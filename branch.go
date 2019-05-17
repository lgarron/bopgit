package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/logrusorgru/aurora"
)

func currentBranch() Branch {
	return NewBranch(runGitCommand("rev-parse", "--abbrev-ref", "HEAD"))
}

func doesBranchContain(branch Branch, ref string) bool {
	return isGitCommandExitCodeZero("merge-base", "--is-ancestor", ref, branch.Name)
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
	runGitCommand("rebase", "--onto", newbase, upstream, root.Name)
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
