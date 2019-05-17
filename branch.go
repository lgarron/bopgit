package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/logrusorgru/aurora"
)

func currentBranch() Branch {
	return NewBranch(getGitValue("rev-parse", "--abbrev-ref", "HEAD"))
}

func doesBranchContain(branch Branch, commit Commit) bool {
	return isGitCommandExitCodeZero("merge-base", "--is-ancestor", commit.Hash, branch.Name)
}

func branchMustContain(branch Branch, commit Commit) {
	if !doesBranchContain(branch, commit) {
		fmt.Fprintf(os.Stderr, "Branch %s does not contain expected ref: %s\n",
			aurora.Bold(branch),
			aurora.Bold(commit),
		)
		os.Exit(1)
	}
}

func rebaseOnto(newbase Commit, upstream Commit, root Branch) {
	runGitCommand("rebase", "--onto", newbase.Hash, upstream.Hash, root.Name)
}

func numCommitsAhead(branch Ref, comparison Ref) int {
	s := getGitValue("rev-list", "--left-only", "--count", fmt.Sprintf(
		"%s...%s",
		branch.ID(),
		comparison.ID(),
	))
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	return i
}
