package main

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
)

func hash(ref string) string {
	return runGitCommand("show-ref", "--heads", "-s", ref)
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

/******** Ref *********/

type Ref struct {
	ID string
}

func (r Ref) String() string {
	return fmt.Sprintf("⇢ %s", r.ID)
}

func NewRef(refID string) Ref {
	// TODO: Check for existence of Ref.
	return Ref{
		ID: refID,
	}
}

func (r Ref) Commit() Commit {
	return NewCommit(hash(r.ID))
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
