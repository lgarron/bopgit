package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/xlab/treeprint"
)

func maybeNumCommitsLeftAheadStr(left Ref, right Ref) string {
	leftAhead, err := maybeNumCommitsLeftAhead(execOptions{timeout: 100 * time.Millisecond}, left, right)
	aheadStr := strconv.Itoa(leftAhead)
	if err != nil {
		aheadStr = "???"
	}
	return aheadStr
}

func maybeNumCommitsDiffStr(left Ref, right Ref) (string, string) {
	leftAhead, rightAhead, err := maybeNumCommitsDiff(execOptions{timeout: 100 * time.Millisecond}, left, right)
	if err != nil {
		return "???", "???"
	}
	return strconv.Itoa(leftAhead), strconv.Itoa(rightAhead)
}

func maybeNumCommitsDiffStrColored(left Ref, right Ref) string {
	leftAhead, rightAhead := maybeNumCommitsDiffStr(left, right)
	return fmt.Sprintf("%s/%s",
		aurora.Sprintf(aurora.Red("+%s"), leftAhead),
		aurora.Sprintf(aurora.Green("-%s"), rightAhead),
	)
}

type branchInfo struct {
	branch     Branch
	baseBranch Branch
	info       string
}

func getBranchInfo(branch Branch) branchInfo {
	baseBranch := getSymBase(branch)
	if !showDistancesInTree {
		return branchInfo{
			branch:     branch,
			baseBranch: baseBranch,
			info:       branch.Name,
		}
	}

	suffix := ""
	if showUpstreamDistances {
		upstream, err := mabyeGetBranch(execOptions{},
			"rev-parse", "--abbrev-ref",
			fmt.Sprintf("%s@{upstream}", branch.Name),
		)
		if err == nil {
			suffix = fmt.Sprintf(" | local: %s",
				maybeNumCommitsDiffStrColored(upstream, branch),
			)
		}
	}

	commitPluralized := "commits"
	commitsOnBranchStr := maybeNumCommitsLeftAheadStr(branch, getLatestBaseCommit(branch))
	if commitsOnBranchStr == "1" {
		commitPluralized = "commit"
	}

	info := fmt.Sprintf("%s | %s | %s %s%s",
		maybeNumCommitsDiffStrColored(baseBranch, branch),
		aurora.Bold(branch.Name),
		commitsOnBranchStr,
		commitPluralized,
		suffix,
	)

	return branchInfo{
		branch:     branch,
		baseBranch: baseBranch,
		info:       info,
	}
}

type branchInfoFuture chan branchInfo
type branchInfoLookup struct {
	futures map[string]branchInfoFuture
}

func newBranchInfoLookup(branches []Branch) branchInfoLookup {
	lookup := branchInfoLookup{
		futures: map[string]branchInfoFuture{},
	}

	for _, branch := range branches {
		lookup.futures[branch.Name] = make(branchInfoFuture)
	}
	return lookup
}

func (l branchInfoLookup) Set(branch Branch, info branchInfo) {
	l.futures[branch.Name] <- info
}

func (l branchInfoLookup) Get(branch Branch) (branchInfo, bool) {
	future, present := l.futures[branch.Name]
	if !present {
		return branchInfo{}, false
	}
	info := <-future
	go func(info branchInfo) {
		future <- info
	}(info)
	return info, true
}

// From http://www.golangpatterns.info/concurrency/parallel-for-loop#TOC-Usage
func calculateBranchInfo(branches []Branch) branchInfoLookup {
	lookup := newBranchInfoLookup(branches)

	// Throttle the amount of parallelism?
	for _, branch := range branches {
		go func(branch Branch) {
			lookup.Set(branch, getBranchInfo(branch))
		}(branch)
	}

	return lookup
}

func ensureInTree(t treeprint.Tree, nodeMemo map[string]treeprint.Tree, lookup branchInfoLookup, branch Branch) treeprint.Tree {
	node := nodeMemo[branch.Name]
	if node != nil {
		return node
	}
	info, present := lookup.Get(branch)
	if !present {
		// New top-level
		newNode := t.AddBranch(branch.Name)
		nodeMemo[branch.Name] = newNode
		return newNode
	}
	parentNode := ensureInTree(t, nodeMemo, lookup, info.baseBranch)
	newNode := parentNode.AddBranch(info.info)
	nodeMemo[branch.Name] = newNode
	return newNode
}

func list() {
	branches := bopgitBranches()
	branchInfoLookup := calculateBranchInfo(branches)

	t := treeprint.New()
	nodeMemo := map[string]treeprint.Tree{}

	for _, branch := range bopgitBranches() {
		ensureInTree(t, nodeMemo, branchInfoLookup, branch)
	}

	fmt.Printf("%s\n", t)
}
