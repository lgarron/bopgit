package main

import (
	"flag"
	"fmt"

	"github.com/logrusorgru/aurora"
)

var debug = false

func showHelp() {
	fmt.Println(`Usage:
  set [optional arguments] baseRef
  set [optional arguments] branch baseRef
  set [optional arguments] branch baseRef latest-base-commit
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
		setCmd()
	}
}

func setCmd() {
	if flag.NArg() == 2 {
		set(currentBranch(), flag.Arg(1), nil)
	} else if flag.NArg() == 3 {
		set(flag.Arg(1), flag.Arg(2), nil)
	} else if flag.NArg() == 4 {
		arg3 := flag.Arg(3)
		set(flag.Arg(1), flag.Arg(2), &arg3)
	} else {
		showHelp()
	}
}

func set(branch gitBranch, baseBranch gitBranch, maybeLatestBaseCommit *string) {
	branchMustExist(branch)
	branchMustExist(baseBranch)

	fmt.Printf("Setting the base branch for %s to %s\n",
		aurora.Bold(branch),
		aurora.Bold(baseBranch),
	)

	setSymBase(branch, baseBranch, "bopgit set")

	var latestBaseCommit string
	if maybeLatestBaseCommit == nil {
		latestBaseCommit = hash(baseBranch)
		fmt.Printf("Calculated latest base commit: %s\n",
			aurora.Bold(latestBaseCommit),
		)
	} else {
		latestBaseCommit = *maybeLatestBaseCommit
		fmt.Printf("Using the latest base commit provided: %s\n",
			aurora.Bold(latestBaseCommit),
		)
	}
	branchMustContain(branch, latestBaseCommit)
	setLatestBaseCommit(branch, latestBaseCommit, "bopgit set")
}

func updateCmd() {
	if flag.NArg() == 1 {
		update(currentBranch())
	} else if flag.NArg() == 2 {
		update(flag.Arg(1))
	} else {
		showHelp()
	}
}

func update(branch gitBranch) {
	branchMustExist(branch)

	fmt.Printf("Updating branch %s\n",
		aurora.Bold(branch),
	)

	baseBranch := getSymBase(branch)
	newLatestBaseCommit := hash(baseBranch)
	oldLatestBaseCommit := getLatestBaseCommit(branch)
	runGitCommand("rebase", "--into", newLatestBaseCommit, oldLatestBaseCommit, branch)

	setLatestBaseCommit(branch, newLatestBaseCommit, "bopgit update")
}
