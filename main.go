package main

import (
	"flag"
	"fmt"

	"github.com/logrusorgru/aurora"
)

var debug = false

func showHelp() {
	fmt.Println(`Usage:
  set [optional arguments] base
  set [optional arguments] base latest-base-commit
  set [optional arguments] base latest-base-commit branch
  update [optional arguments]
  update [optional arguments] branch

  Optional arguments:
    --debug`)
}

func mustHaveMinNArgs(n int) {
	if flag.NArg() < n {
		showHelp()
	}
}

func main() {
	debugPtr := flag.Bool("debug", false, "debug")
	flag.Parse()
	mustHaveMinNArgs(1)

	debug = *debugPtr

	switch flag.Arg(0) {
	case "help":
		showHelp()
	case "set":
		setCmd()
	case "update":
		updateCmd()
	}
}

func setCmd() {
	if flag.NArg() < 2 || flag.NArg() > 4 {
		showHelp()
	}

	baseBranch := currentBranch()
	branchMustExist(baseBranch)

	var branch gitBranch
	if flag.NArg() > 3 {
		branch = flag.Arg(3)
	} else {
		branch = currentBranch()
	}
	branchMustExist(branch)

	fmt.Printf("Setting the base branch for %s to %s\n",
		aurora.Bold(branch),
		aurora.Bold(baseBranch),
	)

	var latestBaseCommit string
	if flag.NArg() > 2 {
		latestBaseCommit = flag.Arg(2)
		fmt.Printf("Using the latest base commit provided: %s\n",
			aurora.Bold(latestBaseCommit),
		)
	} else {
		latestBaseCommit = hash(baseBranch)
		fmt.Printf("Calculated latest base commit: %s\n",
			aurora.Bold(latestBaseCommit),
		)
	}

	branchMustContain(branch, latestBaseCommit)
	set(baseBranch, latestBaseCommit, branch)
}

func set(baseBranch gitBranch, latestBaseCommit string, branch gitBranch) {
	setSymBase(branch, baseBranch, "bopgit set")
	setLatestBaseCommit(branch, latestBaseCommit, "bopgit set")
}

func updateCmd() {
	if flag.NArg() < 2 || flag.NArg() > 3 {
		showHelp()
	}

	var branch gitBranch
	if flag.NArg() > 2 {
		branch = flag.Arg(2)
	} else {
		branch = currentBranch()
	}
	branchMustExist(branch)

	fmt.Printf("Updating branch %s\n",
		aurora.Bold(branch),
	)

	update(branch)
}

func update(branch gitBranch) {
	baseBranch := getSymBase(branch)
	newLatestBaseCommit := hash(baseBranch)
	oldLatestBaseCommit := getLatestBaseCommit(branch)
	rebaseOnto(newLatestBaseCommit, oldLatestBaseCommit, branch)
	setLatestBaseCommit(branch, newLatestBaseCommit, "bopgit update")
}
