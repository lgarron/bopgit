package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"
)

const defaultTimeout = 5 * time.Second

type execOptions struct {
	timeout time.Duration
}

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

func maybeGetGitValue(options execOptions, args ...string) (string, error) {
	cmd := gitExecCommand(args...)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	if err := cmd.Start(); err != nil {
		return "", err
	}
	timeout := options.timeout
	if timeout == 0 {
		timeout = defaultTimeout
	}
	timer := time.AfterFunc(timeout, func() {
		cmd.Process.Kill()
	})
	if err := cmd.Wait(); err != nil {
		return "", err
	}
	timer.Stop()

	return strings.TrimSuffix(outb.String(), "\n"), nil
}

func getGitValue(args ...string) string {
	output, err := maybeGetGitValue(execOptions{}, args...)
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
