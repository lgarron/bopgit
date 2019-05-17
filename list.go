package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/xlab/treeprint"
)

func maybeNumCommitsAheadStr(branch Ref, comparison Ref) string {
	ahead, err := maybeNumCommitsAhead(execOptions{timeout: 200 * time.Millisecond}, branch, comparison)
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
	commitPluralized := "commits"
	commitsOnBranchStr := maybeNumCommitsAheadStr(branch, getLatestBaseCommit(branch))
	if commitsOnBranchStr == "1" {
		commitPluralized = "commit"
	}
	if showDistancesInTree {
		metaText := fmt.Sprintf("%s [%s / %s] (%s %s)",
			aurora.Bold(branch.Name),
			aurora.Sprintf(aurora.Red("-%s"), maybeNumCommitsAheadStr(baseBranch, branch)),
			aurora.Sprintf(aurora.Green("+%s"), maybeNumCommitsAheadStr(branch, baseBranch)),
			commitsOnBranchStr,
			commitPluralized,
		)
		// newNode = parentNode.AddMetaBranch(metaText, )
		newNode = parentNode.AddBranch(metaText)
	} else {
		newNode = parentNode.AddBranch(branch.Name)
	}
	nodeMemo[branch.Name] = newNode
	return newNode
}

func list() {
	t := treeprint.New()
	nodeMemo := map[string]treeprint.Tree{}
	for _, branch := range bopgitBranches() {
		ensureInTree(t, nodeMemo, branch)
	}

	fmt.Printf("%s\n", t)
}
