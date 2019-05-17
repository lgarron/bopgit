package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
)

func showHelpAndExit() {
	fmt.Println(`Usage:
    set [baseRef]
    set [branch] [baseRef]`)
	os.Exit(0)
}

func mustHaveMinNArgs(n int) {
	if flag.NArg() < n {
		showHelpAndExit()
	}
}

func main() {
	flag.Parse()
	mustHaveMinNArgs(1)

	switch flag.Arg(0) {
	case "help":
		showHelpAndExit()
	case "set":
		setCmd()
	case "update":
		setCmd()
	}
}

func setCmd() {
	if flag.NArg() == 2 {
		set(currentBranch(), flag.Arg(1))
	} else if flag.NArg() == 3 {
		set(flag.Arg(1), flag.Arg(2))
	} else {
		showHelpAndExit()
	}
}

func set(branch gitBranch, baseBranch gitBranch) {
	branchMustExist(branch)
	branchMustExist(baseBranch)

	fmt.Printf("Setting the base branch for %s to %s\n",
		aurora.Bold(branch),
		aurora.Bold(baseBranch),
	)

	setSymBase(branch, baseBranch, "bopgit set")

	latestBaseCommit := hash(baseBranch)

	fmt.Printf("Latest base commit is %s\n",
		aurora.Bold(latestBaseCommit),
	)
	setLatestBaseCommit(branch, latestBaseCommit, "bopgit set")
}

func updateCmd() {
	if flag.NArg() == 1 {
		update(currentBranch())
	} else if flag.NArg() == 2 {
		update(flag.Arg(1))
	} else {
		showHelpAndExit()
	}
}

func update(branch gitBranch) {
	branchMustExist(branch)

	fmt.Printf("Updating branch %s\n",
		aurora.Bold(branch),
	)

	// setSymBase(branch, baseBranch, "bopgit set")
	// setLatestBaseCommit(branch, baseBranch, "bopgit set")
}
