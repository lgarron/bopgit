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

func forget(branch Branch) {
	oldBaseBranch := getSymBase(branch)
	oldLatestBaseCommit := getLatestBaseCommit(branch)

	clearSymBase(branch, "bopgit forget")
	clearLatestBaseCommit(branch, "bopgit forget")

	fmt.Println()
	fmt.Printf(aurora.Sprintf(aurora.Green("bopgit has forgotten %s\n"),
		aurora.Bold(branch),
	))

	fmt.Printf(`To restore the previous state, run:
  git checkout %s
  bopgit track %s %s
`,
		aurora.Bold(branch.Name),
		aurora.Bold(oldBaseBranch.Name),
		aurora.Bold(oldLatestBaseCommit.Hash),
	)
}
