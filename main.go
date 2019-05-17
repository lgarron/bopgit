package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/logrusorgru/aurora"
)

func showHelp() {
	fmt.Println(`Usage:
    set [baseRef]
    set [branch] [baseRef]`)
	os.Exit(0)
}

func mustHaveMinNArgs(n int) {
	if flag.NArg() < n {
		showHelp()
	}
}

func main() {
	flag.Parse()
	mustHaveMinNArgs(1)

	switch flag.Arg(0) {
	case "help":
		showHelp()
	case "set":
		setCmd()
	}
}

func runGitCommand(args ...string) string {
	output, err := exec.Command("git", args...).Output()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return strings.TrimSuffix(string(output), "\n")
}

func currentBranch() string {
	return runGitCommand("rev-parse", "--abbrev-ref", "HEAD")
}

func bopgitSymBaseRefName(branch string) string {
	return fmt.Sprintf("refs/bopgit/sym-base/%s", branch)
}

func bopgitCurrentBaseRefName(branch string) string {
	return fmt.Sprintf("refs/bopgit/current-base/%s", branch)
}

func setCmd() {
	if flag.NArg() == 2 {
		set(currentBranch(), flag.Arg(1))
	} else if flag.NArg() == 3 {
		set(flag.Arg(1), flag.Arg(2))
	} else {
		showHelp()
	}
}

func set(branch string, baseRef string) {
	fmt.Printf("Setting the base branch for %s to %s\n",
		aurora.Bold(branch),
		aurora.Bold(baseRef),
	)

	runGitCommand("symbolic-ref", "-m", "bopgit set", bopgitSymBaseRefName(branch), baseRef)
	runGitCommand("update-ref", bopgitCurrentBaseRefName(branch), baseRef)
	// git symbolic-ref [-m <reason>] <name> <ref>
	// fmt.Println(currentBranch(), baseRef)
}
