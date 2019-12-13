package main

import (
	"fmt"
	"log"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

func main() {
	out, err := execOut("git", "status", "--porcelain")
	if err != nil {
		log.Fatal(err)
	}
	repo, err := getRepoPath()
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(out), "\n")
	for i, _ := range lines {
		line := strings.Fields(lines[i])
		if len(line) != 2 || !strings.HasSuffix(line[1], ".go") {
			fmt.Println("Skipping", strconv.Quote(lines[i]))
			continue
		}
		path := path.Join(repo, line[1])
		err := prettify(path)
		if err != nil {
			fmt.Printf("Failed to prettify %s for %s\n", strconv.Quote(path), err.Error())
			continue
		}
		fmt.Println("Prettified ", strconv.Quote(path))
	}
}

func getRepoPath() (string, error) {
	out, err := execOut("git", "rev-parse", "--show-toplevel")
	if err != nil {
		return "", nil
	}
	return out, nil
}

func prettify(file string) error {
	return execErr("gofmt", "-w", file)
}

func execErr(args ...string) error {
	_, err := execOut(args...)
	return err
}

func execOut(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("Not enough arguments")
	}
	out, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
