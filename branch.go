package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

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

func maybeNumCommitsAhead(options execOptions, branch Ref, comparison Ref) (int, error) {
	s, err := maybeGetGitValue(options, "rev-list", "--left-only", "--count", fmt.Sprintf(
		"%s...%s",
		branch.ID(),
		comparison.ID(),
	))
	if err != nil {
		return 0, err
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return i, nil
}

func numCommitsAhead(branch Ref, comparison Ref) int {
	i, err := maybeNumCommitsAhead(execOptions{}, branch, comparison)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return i
}

func bopgitBranches() []Branch {
	// TODO: Look up the `sym-base` refs instead?
	// (Would probably require fixing https://github.com/lgarron/bopgit/issues/3)
	branchesStr := getGitValue("for-each-ref", "--format=%(refname:short)", "refs/bopgit/latest-base-commit")
	branchRefs := strings.Split(branchesStr, "\n")

	branches := []Branch{}
	for _, branchRef := range branchRefs {
		branchName := strings.TrimPrefix(branchRef, "bopgit/latest-base-commit/")
		branches = append(branches, NewBranch(branchName))
	}

	return branches
}
