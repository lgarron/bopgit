package main

import "fmt"

func symBaseRefName(branch gitBranch) string {
	return fmt.Sprintf("refs/bopgit/sym-base/%s", branch)
}

func latestBaseCommitRefName(branch gitBranch) string {
	return fmt.Sprintf("refs/bopgit/latest-base-commit/%s", branch)
}

func setSymBase(branch gitBranch, baseBranch gitBranch, reason string) {
	runGitCommand("symbolic-ref", "-m", reason, symBaseRefName(branch), baseBranch)
}

func setLatestBaseCommit(branch gitBranch, commit string, reason string) {
	runGitCommand("update-ref", "-m", reason, latestBaseCommitRefName(branch), commit)
}

func getSymBase(branch gitBranch) string {
	return runGitCommand("symbolic-ref", symBaseRefName(branch))
}

func getLatestBaseCommit(branch gitBranch) string {
	return runGitCommand("rev-parse", latestBaseCommitRefName(branch))
}