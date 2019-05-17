package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/logrusorgru/aurora"
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
		fmt.Printf(aurora.Sprintf(aurora.Yellow("⚙️  git command: %v\n"), args))
	}
	return exec.Command("git", args...)
}

func maybeGetGitValue(args ...string) (string, error) {
	output, err := gitExecCommand(args...).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(output), "\n"), nil
}

func getGitValue(args ...string) string {
	output, err := maybeGetGitValue(args...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return output
}

func runGitCommand(args ...string) {
	cmd := gitExecCommand(args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func isGitCommandExitCodeZero(args ...string) bool {
	cmd := gitExecCommand(args...)
	err := cmd.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unexpected error while checking exit code of a git command: %s\n", err)
		os.Exit(1)
	}
	err = cmd.Wait()
	if err != nil {
		return false
	}
	return true
}
