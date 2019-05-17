package main

import "fmt"

func symBaseRefName(branch Branch) string {
	return fmt.Sprintf("refs/bopgit/sym-base/%s", branch.name)
}

func latestBaseCommitRefName(branch Branch) string {
	return fmt.Sprintf("refs/bopgit/latest-base-commit/%s", branch.name)
}

func setSymBase(branch Branch, baseBranch Branch, reason string) {
	runGitCommand("symbolic-ref", "-m", reason, symBaseRefName(branch), baseBranch.name)
}

func setLatestBaseCommit(branch Branch, commit string, reason string) {
	runGitCommand("update-ref", "-m", reason, latestBaseCommitRefName(branch), commit)
}

func getSymBase(branch Branch) Branch {
	return newBranch(runGitCommand("symbolic-ref", symBaseRefName(branch)))
}

func getLatestBaseCommit(branch Branch) string {
	return runGitCommand("rev-parse", latestBaseCommitRefName(branch))
}
