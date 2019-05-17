package main

import "os/exec"

type gitBranch = string

func hash(ref string) string {
	return runGitCommand("show-ref", "-d", "-s", ref)
}

func currentBranch() string {
	return runGitCommand("rev-parse", "--abbrev-ref", "HEAD")
}

func doesBranchExist(branch gitBranch) bool {
	cmd := gitExecCommand("rev-parse", "--verify", branch)
	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode() != 0
		}
	}
	return true
}
