package main

import (
	"flag"
	"fmt"
	"log"
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

	baseBranch := NewBranch(flag.Arg(1))

	var branch Branch
	if flag.NArg() > 3 {
		branch = NewBranch(flag.Arg(3))
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
		latestBaseCommit = hash(baseBranch.Name)
		fmt.Printf("Calculated latest base commit: %s\n",
			aurora.Bold(latestBaseCommit),
		)
	}

	branchMustContain(branch, latestBaseCommit)
	set(baseBranch, latestBaseCommit, branch)
}

func set(baseBranch Branch, latestBaseCommit string, branch Branch) {
	setSymBase(branch, baseBranch, "bopgit set")
	setLatestBaseCommit(branch, latestBaseCommit, "bopgit set")
}

func updateCmd() {
	if flag.NArg() < 1 || flag.NArg() > 2 {
		showHelp()
		os.Exit(1)
	}

	var branch Branch
	if flag.NArg() > 1 {
		branch = NewBranch(flag.Arg(1))
	} else {
		branch = currentBranch()
	}

	fmt.Printf("Updating branch %s\n",
		aurora.Bold(branch),
	)

	update(branch)
}

func update(branch Branch) {
	baseBranch := getSymBase(branch)
	newLatestBaseCommit := hash(baseBranch.Name)
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

	var branch Branch
	if flag.NArg() > 1 {
		branch = NewBranch(flag.Arg(1))
	} else {
		branch = currentBranch()
	}

	fmt.Printf("Info for branch %s\n",
		aurora.Bold(branch),
	)

	info(branch)
}

func info(branch Branch) {
	// TODO: Calculate if branch is tracked by `bopgit`.

	baseBranch := getSymBase(branch)
	fmt.Printf("Base branch: %s\n",
		aurora.Bold(baseBranch),
	)

	latestBaseCommit := getLatestBaseCommit(branch)
	fmt.Printf("Latest base commit: %s\n",
		aurora.Bold(getLatestBaseCommit(branch)),
	)

	if !doesBranchContain(branch, latestBaseCommit) {
		log.Fatal("The branch doesn't contain that `bopgit` believes to be its latest base commit!")
		os.Exit(1)
	}

	fmt.Println()

	// TODO: avoid assuming a linear history?
	fmt.Printf("%d commits to %s since the its base commit.\n",
		numCommitsAhead(branch.Name, latestBaseCommit),
		aurora.Bold(branch),
	)

	fmt.Printf("%d commits in %s that %s doesn't have.\n",
		numCommitsAhead(branch.Name, baseBranch.Name),
		aurora.Bold(branch),
		aurora.Bold(baseBranch),
	)

	fmt.Printf("%d commits in %s that %s doesn't have.\n",
		numCommitsAhead(baseBranch.Name, branch.Name),
		aurora.Bold(baseBranch),
		aurora.Bold(branch),
	)
}
