package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func showHelp() {
	flag.Usage()
	os.Exit(0)
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		showHelp()
	}
	switch flag.Arg(0) {
	case "help":
		showHelp()
	case "set":
		set(flag.Arg(1))
	}
}

func runGitCommand(args ...string) {
	out, err := exec.Command("git", args...).Output()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func currentBranch(): string {
  runGitCommand("rev-parse", "--abbrev-ref", "HEAD")
}

func set(baseRef string) {
	fmt.Println(currentBranch(), baseRef)
}
