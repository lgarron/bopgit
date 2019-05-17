package main

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
)

func hash(ref string) string {
	return runGitCommand("show-ref", "--heads", "-s", ref)
}

type Ref interface {
	ID() string
}

/******** Commit *********/

type Commit struct {
	Hash string
}

func (c Commit) String() string {
	return fmt.Sprintf("⌥ %s", c.Hash)
}

func NewCommit(hashStr string) Commit {
	// TODO: Check for existence of commit.
	return Commit{
		Hash: hashStr,
	}
}

func (c Commit) ID() string {
	return c.Hash
}

func (c Commit) Equals(c2 Commit) bool {
	return c.Hash == c2.Hash
}

/******** ArbitraryRef *********/

type ArbitraryRef struct {
	refID string
}

func (r ArbitraryRef) String() string {
	return fmt.Sprintf("⇢ %s", r.refID)
}

func NewArbitraryRef(refID string) ArbitraryRef {
	// TODO: Check for existence of ArbitraryRef.
	return ArbitraryRef{
		refID: refID,
	}
}

func (r ArbitraryRef) ID() string {
	return r.refID
}

/******** Branch *********/

type Branch struct {
	Name string
}

func (b Branch) String() string {
	return fmt.Sprintf("⌥ %s", b.Name)
}

func doesBranchNameExist(branchName string) bool {
	return isGitCommandExitCodeZero("rev-parse", "--verify", branchName)
}

func branchNameMustExist(branchName string) {
	if !doesBranchNameExist(branchName) {
		fmt.Printf("Branch does not exist: %s",
			aurora.Bold(branchName),
		)
		showHelp()
		os.Exit(1)
	}
}

func NewBranch(branchName string) Branch {
	branchNameMustExist(branchName)
	return Branch{
		Name: branchName,
	}
}

func (b Branch) Commit() Commit {
	return NewCommit(hash(b.Name))
}

func (b Branch) ID() string {
	return b.Name
}
