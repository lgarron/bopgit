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

func maybeNumCommitsLeftAhead(options execOptions, left, right Ref) (int, error) {
	s, err := maybeGetGitValue(options, "rev-list", "--left-only", "--count", fmt.Sprintf(
		"%s...%s",
		left.ID(),
		right.ID(),
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

func maybeNumCommitsDiff(options execOptions, left, right Ref) (int, int, error) {
	s, err := maybeGetGitValue(options, "rev-list", "--left-right", "--count", fmt.Sprintf(
		"%s...%s",
		left.ID(),
		right.ID(),
	))
	if err != nil {
		return 0, 0, err
	}
	diff := strings.Fields(s)
	leftAhead, err := strconv.Atoi(diff[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	rightAhead, err := strconv.Atoi(diff[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return leftAhead, rightAhead, nil
}

func numCommitsLeftAhead(left, right Ref) int {
	i, err := maybeNumCommitsLeftAhead(execOptions{}, left, right)
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
