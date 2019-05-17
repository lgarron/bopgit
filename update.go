package main

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
)

func update(branch Branch) {
	baseBranch, err := mabyeGetSymBase(branch)
	if err != nil {
		fmt.Println("bopgit is not tracking this branch.")
		os.Exit(1)
	}
	newLatestBaseCommit := baseBranch.Commit()
	oldLatestBaseCommit := getLatestBaseCommit(branch)

	if newLatestBaseCommit.Equals(oldLatestBaseCommit) {
		fmt.Printf("Base commit %s is up to date.\n",
			aurora.Bold(oldLatestBaseCommit),
		)
		os.Exit(0)
	}

	fmt.Println()

	fmt.Printf("Old base commit: %s\n",
		aurora.Bold(oldLatestBaseCommit),
	)

	fmt.Printf("New base commit: %s\n",
		aurora.Bold(newLatestBaseCommit),
	)

	fmt.Println()

	oldHeadCommit := branch.Commit()
	fmt.Printf("Old HEAD: %s\n",
		aurora.Bold(oldHeadCommit),
	)

	// TODO: track backup ref.
	rebaseOnto(newLatestBaseCommit, oldLatestBaseCommit, branch)
	setLatestBaseCommit(branch, newLatestBaseCommit, "bopgit update")

	fmt.Printf("New HEAD: %s\n",
		aurora.Bold(branch.Commit()),
	)

	fmt.Println()
	fmt.Printf(aurora.Sprintf(aurora.Green("✅ Updated %s\n"),
		aurora.Green(aurora.Bold(branch)),
	))

	fmt.Println()

	fmt.Printf(`To restore to the previous state, run:
  git checkout %s
  git reset --hard %s
  bopgit track %s %s %s
`,
		aurora.Bold(branch.Name),
		aurora.Bold(oldHeadCommit.Hash),
		aurora.Bold(baseBranch.Name),
		aurora.Bold(oldLatestBaseCommit.Hash),
		aurora.Bold(branch.Name),
	)
}
