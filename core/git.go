package core

import (
	"os/exec"
	"strconv"
	"strings"
)

func runGit(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	out, cmd.Err = cmd.Output()
	if cmd.Err != nil {
		return "", cmd.Err
	}
	return strings.TrimSpace(string(out)), nil
}

func IsGitRepo() bool {
	_, err := runGit("rev-parse", "--is-inside-work-tree")
	return err == nil
}

func GetMetadata() (MetaData, error) {
	repo, _ := runGit("remote", "get-url", "origin")
	branch, _ := runGit("branch", "--show-current")
	lastCommit, _ := runGit("log", "-1", "--format=%h")
	lastCommitMsg, _ := runGit("log", "-1", "--format=%s")
	lastCommitDate, _ := runGit("log", "-1", "--format=%ci")
	commitsRaw, _ := runGit("rev-list", "--count", "HEAD")
	commits, _ := strconv.Atoi(commitsRaw)
	contributors := getContributors()

	return MetaData{
		Version:           "1.0.0",
		Repository:        repo,
		Branch:            branch,
		Commits:           commits,
		LastCommit:        lastCommit,
		LastCommitMessage: lastCommitMsg,
		LastCommitDate:    lastCommitDate,
		Contributors:      contributors,
	}, nil
}

func getContributors() []string {
	raw, err := runGit("log", "--format=%an")
	if err != nil {
		return []string{}
	}

	lines := strings.Split(raw, "\n")
	seen := map[string]bool{}
	var contributors []string
	for _, name := range lines {
		if name != "" && !seen[name] {
			seen[name] = true
			contributors = append(contributors, name)
		}
	}
	return contributors
}
