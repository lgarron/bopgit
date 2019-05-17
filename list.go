package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/xlab/treeprint"
)

func maybeNumCommitsAheadStr(branch Ref, comparison Ref) string {
	ahead, err := maybeNumCommitsAhead(execOptions{timeout: 50 * time.Millisecond}, branch, comparison)
	aheadStr := strconv.Itoa(ahead)
	if err != nil {
		aheadStr = "???"
	}
	return aheadStr
}

func branchInfo(branch Branch) string {
	if !showDistancesInTree {
		return branch.Name
	}

	baseBranch := getSymBase(branch)

	suffix := ""
	if showUpstreamDistances {
		upstream, err := mabyeGetBranch(execOptions{},
			"rev-parse", "--abbrev-ref",
			fmt.Sprintf("%s@{upstream}", branch.Name),
		)
		if err == nil {
			suffix = fmt.Sprintf(" | upstream: %s/%s",
				aurora.Sprintf(aurora.Red("+%s"), maybeNumCommitsAheadStr(branch, upstream)),
				aurora.Sprintf(aurora.Green("-%s"), maybeNumCommitsAheadStr(upstream, branch)),
			)
		}
	}

	commitPluralized := "commits"
	commitsOnBranchStr := maybeNumCommitsAheadStr(branch, getLatestBaseCommit(branch))
	if commitsOnBranchStr == "1" {
		commitPluralized = "commit"
	}

	return fmt.Sprintf("%s/%s | %s | %s %s%s",
		aurora.Sprintf(aurora.Red("-%s"), maybeNumCommitsAheadStr(baseBranch, branch)),
		aurora.Sprintf(aurora.Green("+%s"), maybeNumCommitsAheadStr(branch, baseBranch)),
		aurora.Bold(branch.Name),
		commitsOnBranchStr,
		commitPluralized,
		suffix,
	)
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
	newNode := parentNode.AddBranch(branchInfo(branch))
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
