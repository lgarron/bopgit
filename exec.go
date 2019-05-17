package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func Log(args ...interface{}) {
	for i, a := range args {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(a)
	}
	fmt.Print("\n")
}

func gitExecCommand(args ...string) *exec.Cmd {
	if debug {
		fmt.Printf("git command: %v\n", args)
	}
	return exec.Command("git", args...)
}

func runGitCommand(args ...string) string {
	output, err := gitExecCommand(args...).Output()
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSuffix(string(output), "\n")
}

func isGitCommandExitCodeZero(args ...string) bool {
	cmd := gitExecCommand(args...)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Wait()
	if err != nil {
		return false
	}
	return true
}
