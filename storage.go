package main

import "fmt"

func symBaseRefName(branch Branch) string {
	return fmt.Sprintf("refs/bopgit/sym-base/%s", branch.Name)
}

func latestBaseCommitRefName(branch Branch) string {
	return fmt.Sprintf("refs/bopgit/latest-base-commit/%s", branch.Name)
}

func setSymBase(branch Branch, baseBranch Branch, reason string) {
	getGitValue("symbolic-ref", "-m", reason, symBaseRefName(branch), baseBranch.Name)
}

func setLatestBaseCommit(branch Branch, commit Commit, reason string) {
	getGitValue("update-ref", "-m", reason, latestBaseCommitRefName(branch), commit.Hash)
}

func mabyeGetSymBase(branch Branch) (Branch, error) {
	return mabyeGetBranch(execOptions{}, "symbolic-ref", symBaseRefName(branch))
}

func getSymBase(branch Branch) Branch {
	return NewBranch(getGitValue("symbolic-ref", symBaseRefName(branch)))
}

func getLatestBaseCommit(branch Branch) Commit {
	return NewCommit(getGitValue("rev-parse", latestBaseCommitRefName(branch)))
}
