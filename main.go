package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/xlab/treeprint"
)

var debug = false
var showDistancesInTree = true

func showHelp() {
	fmt.Println(`Usage:
  track base-branch
  track base-branch latest-base-commit
  info
  info branch
  update
  tree

  Optional arguments (before positional):
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
	case "track":
		trackCmd()
	case "update":
		updateCmd()
	case "info":
		infoCmd()
	case "tree":
		treeCmd()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", flag.Arg(0))
		os.Exit(1)
	}
}

func trackCmd() {
	if flag.NArg() < 2 || flag.NArg() > 3 {
		showHelp()
		os.Exit(1)
	}

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

func track(baseBranch Branch, latestBaseCommit Commit, branch Branch) {
	setSymBase(branch, baseBranch, "bopgit track")
	setLatestBaseCommit(branch, latestBaseCommit, "bopgit track")

	fmt.Println()
	fmt.Printf(aurora.Sprintf(aurora.Green("✅ %s is now tracking %s\n"),
		aurora.Green(aurora.Bold(branch)),
		aurora.Green(aurora.Bold(baseBranch)),
	))
}

func updateCmd() {
	if flag.NArg() < 1 || flag.NArg() > 1 {
		showHelp()
		os.Exit(1)
	}

	branch := currentBranch()

	fmt.Printf("Updating branch %s\n",
		aurora.Bold(branch),
	)

	update(branch)
}

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

	fmt.Printf("ℹ️  Info for branch %s\n",
		aurora.Bold(branch),
	)

	info(branch)
}

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
		numCommitsAhead(branch, latestBaseCommit),
		aurora.Bold(branch),
	)

	fmt.Printf("%d commits in %s that %s doesn't have.\n",
		numCommitsAhead(branch, baseBranch),
		aurora.Bold(branch),
		aurora.Bold(baseBranch),
	)

	fmt.Printf("%d commits in %s that %s doesn't have.\n",
		numCommitsAhead(baseBranch, branch),
		aurora.Bold(baseBranch),
		aurora.Bold(branch),
	)
}

func treeCmd() {
	if flag.NArg() < 1 || flag.NArg() > 1 {
		showHelp()
		os.Exit(1)
	}

	tree()
}

func maybeNumCommitsAheadStr(branch Ref, comparison Ref) string {
	ahead, err := maybeNumCommitsAhead(execOptions{timeout: 100 * time.Millisecond}, branch, comparison)
	aheadStr := strconv.Itoa(ahead)
	if err != nil {
		aheadStr = "???"
	}
	return aheadStr
}

func ensureInTree(t treeprint.Tree, nodeMemo map[string]treeprint.Tree, branch Branch) treeprint.Tree {
	node := nodeMemo[branch.Name]
	if node != nil {
		return node
	}
	baseBranch, err := mabyeGetSymBase(branch)
	if err != nil {
		// New top-level
		newNode := t.AddBranch(branch.Name)
		nodeMemo[branch.Name] = newNode
		return newNode
	}
	parentNode := ensureInTree(t, nodeMemo, baseBranch)

	var newNode treeprint.Tree
	if showDistancesInTree {
		metaText := fmt.Sprintf("-%s, +%s",
			maybeNumCommitsAheadStr(baseBranch, branch),
			maybeNumCommitsAheadStr(branch, baseBranch),
		)
		newNode = parentNode.AddMetaBranch(metaText, branch.Name)
	} else {
		newNode = parentNode.AddBranch(branch.Name)
	}
	nodeMemo[branch.Name] = newNode
	return newNode
}

func tree() {
	t := treeprint.New()
	nodeMemo := map[string]treeprint.Tree{}
	for _, branch := range bopgitBranches() {
		ensureInTree(t, nodeMemo, branch)
	}

	fmt.Printf("%s\n", t)
}
