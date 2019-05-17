package main

import "fmt"

func symBaseRefName(branch Branch) string {
	return fmt.Sprintf("refs/bopgit/sym-base/%s", branch.Name)
}

func latestBaseCommitRefName(branch Branch) string {
	return fmt.Sprintf("refs/bopgit/latest-base-commit/%s", branch.Name)
}

func setSymBase(branch Branch, baseBranch Branch, reason string) {
	runGitCommand("symbolic-ref", "-m", reason, symBaseRefName(branch), baseBranch.Name)
}

func setLatestBaseCommit(branch Branch, commit string, reason string) {
	runGitCommand("update-ref", "-m", reason, latestBaseCommitRefName(branch), commit)
}

func getSymBase(branch Branch) Branch {
	return NewBranch(runGitCommand("symbolic-ref", symBaseRefName(branch)))
}

func getLatestBaseCommit(branch Branch) string {
	return runGitCommand("rev-parse", latestBaseCommitRefName(branch))
}
