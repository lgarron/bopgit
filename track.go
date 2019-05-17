package main

import (
	"fmt"

	"github.com/logrusorgru/aurora"
)

func track(baseBranch Branch, latestBaseCommit Commit, branch Branch) {
	setSymBase(branch, baseBranch, "bopgit track")
	setLatestBaseCommit(branch, latestBaseCommit, "bopgit track")

	fmt.Println()
	fmt.Printf(aurora.Sprintf(aurora.Green("✅ %s is now tracking %s\n"),
		aurora.Green(aurora.Bold(branch)),
		aurora.Green(aurora.Bold(baseBranch)),
	))
}
