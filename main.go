package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
)

var debug = false
var showDistancesInTree = true
var showUpstreamDistances = true

const globalMaxNArgs = 4

func showHelp() {
	fmt.Println(`Usage:
  track base-branch
  track base-branch latest-base-commit
  forget
  forget branch
  info
  info branch
  update
  list
  base-branch
  base-branch branch
  latest-base-commit
  latest-base-commit branch

  Optional arguments (before positional):
    --debug`)
}

func mustHaveNArgsInRange(min int, max int) {
	if max > globalMaxNArgs {
		fmt.Errorf("Internal argument count error.")
		os.Exit(1)
	}

	if flag.NArg() < min || flag.NArg() > max {
		showHelp()
		os.Exit(1)
	}
}

func main() {
	debugPtr := flag.Bool("debug", false, "debug")
	flag.Parse()
	mustHaveNArgsInRange(1, 4)
	if flag.NArg() < 1 {
		showHelp()
		os.Exit(0)
	}

	debug = *debugPtr

	switch flag.Arg(0) {
	case "help":
		showHelp()
	case "track":
		trackCmd()
	case "forget":
		forgetCmd()
	case "update":
		updateCmd()
	case "info":
		infoCmd()
	case "list":
		listCmd()
	case "branch-base":
		branchBaseCmd()
	case "latest-base-commit":
		latestBaseCommitCmd()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", flag.Arg(0))
		os.Exit(1)
	}
}

// Defaults to current branch.
func branchArg(idx int, mustExist bool) Branch {
	if flag.NArg() > idx {
		return NewBranchWithExistence(flag.Arg(idx), mustExist)
	}

	return currentBranch()
}

func trackCmd() {
	mustHaveNArgsInRange(2, 3)

	baseBranch := NewBranch(flag.Arg(1))

	branch := currentBranch()

	fmt.Printf("Setting branch %s to track %s\n",
		aurora.Bold(branch),
		aurora.Bold(baseBranch),
	)

	var latestBaseCommit Commit
	if flag.NArg() > 2 {
		ref := NewArbitraryRef(flag.Arg(2))
		latestBaseCommit = ref.Commit()
		fmt.Printf("Using the latest base commit from the provided ref: %s\n",
			aurora.Bold(latestBaseCommit),
		)
	} else {
		latestBaseCommit = baseBranch.Commit()
		fmt.Printf("Calculated latest base commit: %s\n",
			aurora.Bold(latestBaseCommit),
		)
	}

	branchMustContain(branch, latestBaseCommit)
	track(baseBranch, latestBaseCommit, branch)
}

func forgetCmd() {
	mustHaveNArgsInRange(1, 2)

	branch := branchArg(1, false)

	fmt.Printf("Forgetting branch %s\n",
		aurora.Bold(branch),
	)

	// TODO; Update refs that refer to this one.
	forget(branch)
}

func updateCmd() {
	mustHaveNArgsInRange(1, 1)

	branch := currentBranch()

	fmt.Printf("Updating branch %s\n",
		aurora.Bold(branch),
	)

	update(branch)
}

func infoCmd() {
	mustHaveNArgsInRange(1, 2)

	var branch = branchArg(1, true)
	fmt.Printf("ℹ️  Info for branch %s\n",
		aurora.Bold(branch),
	)

	info(branch)
}

func listCmd() {
	mustHaveNArgsInRange(1, 1)
	list()
}

func branchBaseCmd() {
	mustHaveNArgsInRange(1, 2)

	branch := branchArg(1, false)

	fmt.Printf("%s\n", getSymBase(branch).Name)
}

func latestBaseCommitCmd() {
	mustHaveNArgsInRange(1, 2)

	branch := branchArg(1, false)

	fmt.Printf("%s\n", getLatestBaseCommit(branch).Hash)
}
