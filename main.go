package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"

	"github.com/Edditoria/grelease/github"
)

const (
	usageUrl      = "github.com/OWNER/REPO"
	usageFetchCmd = "grelease fetch --repo=" + usageUrl + " --write=<FILE_PATH>"
)

func main() {
	fetch := flag.NewFlagSet("fetch", flag.ExitOnError)
	fetchRepo := fetch.String("repo", "", "URL to fetch: "+usageUrl)
	fetchWrite := fetch.String("write", "", "file path for writing a JSON file")

	if len(os.Args) < 2 {
		os.Stderr.WriteString("Usage: " + usageFetchCmd + "\n")
		os.Exit(1)
	}

	switch os.Args[1] {
	case fetch.Name():
		err := fetch.Parse(os.Args[2:])
		if err != nil {
			os.Stderr.WriteString("grelease: argument error: " + err.Error() + "\n")
			os.Exit(1)
		}
		exitCode := handleFetch(fetchRepo, fetchWrite)
		os.Exit(exitCode)
	default:
		os.Stderr.WriteString("grelease: invalid command\nUsage: " + usageFetchCmd + "\n")
		os.Exit(1)
	}
}

// Subcommand "fetch" to consume API call, then write to a file or `stdout`.
// It also handle errors for users.
//
// @return int for exit code.
func handleFetch(repoFlag, writeFlag *string) int {
	splited := strings.Split(*repoFlag, "/")
	if len(splited) != 3 {
		os.Stderr.WriteString("grelease: invalid repository url\nExpect: " + usageUrl + "\n")
		return 1
	}
	if splited[0] != "github.com" {
		os.Stderr.WriteString("grelease: hostname not supported currently\nExpect: " + usageUrl + "\n")
		return 1
	}

	repo := github.Repo{Owner: splited[1], Name: splited[2]}
	// TODO: fetch all releases here:
	releases, _, err := repo.ListReleases(1)
	if err != nil {
		panic(err)
	}
	repo.Releases = releases

	if *writeFlag == "" {
		err = repo.WriteJson(os.Stdout, "", "\t")
		if err != nil {
			os.Stderr.WriteString(err.Error() + "\n")
			return 1
		}
		return 0
	}

	absPath, err := filepath.Abs(*writeFlag)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		return 1
	}
	err = repo.WriteJsonFile(absPath, "", "\t")
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		return 1
	}

	return 0
}
