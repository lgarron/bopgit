package main

import (
	"flag"
	"fmt"
	"os"

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
  info [optional arguments]
  info [optional arguments] branch

  Optional arguments:
    --debug`)
}

func main() {
	debugPtr := flag.Bool("debug", false, "debug")
	flag.Parse()
	if flag.NArg() < 1 {
		showHelp()
		os.Exit(0)
	}

	debug = *debugPtr

	switch flag.Arg(0) {
	case "help":
		showHelp()
	case "set":
		setCmd()
	case "update":
		updateCmd()
	case "info":
		infoCmd()
	}
}

func setCmd() {
	if flag.NArg() < 2 || flag.NArg() > 4 {
		showHelp()
		os.Exit(1)
	}

	baseBranch := newGitBranch(flag.Arg(1))

	var branch gitBranch
	if flag.NArg() > 3 {
		branch = newGitBranch(flag.Arg(3))
	} else {
		branch = currentBranch()
	}

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
		latestBaseCommit = hash(baseBranch.name)
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
	if flag.NArg() < 1 || flag.NArg() > 2 {
		showHelp()
		os.Exit(1)
	}

	var branch gitBranch
	if flag.NArg() > 1 {
		branch = newGitBranch(flag.Arg(1))
	} else {
		branch = currentBranch()
	}

	fmt.Printf("Updating branch %s\n",
		aurora.Bold(branch),
	)

	update(branch)
}

func update(branch gitBranch) {
	baseBranch := getSymBase(branch)
	newLatestBaseCommit := hash(baseBranch.name)
	oldLatestBaseCommit := getLatestBaseCommit(branch)

	fmt.Printf("Updating branch: %s\n",
		aurora.Bold(branch),
	)

	fmt.Printf("Old base commit: %s\n",
		aurora.Bold(oldLatestBaseCommit),
	)

	fmt.Printf("New base commit: %s\n",
		aurora.Bold(newLatestBaseCommit),
	)

	// TODO: Set backup ref.
	rebaseOnto(newLatestBaseCommit, oldLatestBaseCommit, branch)
	setLatestBaseCommit(branch, newLatestBaseCommit, "bopgit update")
}

func infoCmd() {
	if flag.NArg() < 1 || flag.NArg() > 2 {
		showHelp()
		os.Exit(1)
	}

	var branch gitBranch
	if flag.NArg() > 1 {
		branch = newGitBranch(flag.Arg(1))
	} else {
		branch = currentBranch()
	}

	fmt.Printf("Info for branch %s\n",
		aurora.Bold(branch),
	)

	info(branch)
}

func info(branch gitBranch) {
	// TODO: Calculate if branch is tracked by `bopgit`.

	fmt.Printf("Base branch: %s\n",
		aurora.Bold(getSymBase(branch)),
	)

	fmt.Printf("Latest base commit: %s\n",
		aurora.Bold(getLatestBaseCommit(branch)),
	)
}
