package main

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
)

func info(branch Branch) {
	fmt.Println()
	baseBranch, err := mabyeGetSymBase(branch)
	if err != nil {
		fmt.Println("bopgit is not tracking this branch.")
		os.Exit(0)
	}
	fmt.Printf("Base branch: %s\n",
		aurora.Bold(baseBranch),
	)

	latestBaseCommit := getLatestBaseCommit(branch)
	fmt.Printf("Latest base commit: %s\n",
		aurora.Bold(latestBaseCommit),
	)

	if !doesBranchContain(branch, latestBaseCommit) {
		fmt.Fprintf(os.Stderr, "The branch doesn't contain that `bopgit` believes to be its latest base commit!")
		os.Exit(1)
	}

	fmt.Println()

	// TODO: avoid assuming a linear history?
	fmt.Printf("%d commits to %s since its base commit.\n",
		numCommitsLeftAhead(branch, latestBaseCommit),
		aurora.Bold(branch),
	)

	// TODO: Parallelize this with the comparison against the base commit.
	leftAhead, rightAhead := numCommitsDiff(baseBranch, branch)

	fmt.Printf("%d commits in %s that %s doesn't have.\n",
		leftAhead,
		aurora.Bold(branch),
		aurora.Bold(baseBranch),
	)

	fmt.Printf("%d commits in %s that %s doesn't have.\n",
		rightAhead,
		aurora.Bold(baseBranch),
		aurora.Bold(branch),
	)

	// TODO: Print commit distance against upstream.
}
