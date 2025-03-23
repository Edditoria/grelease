package main

import (
	"flag"
	"os"
	"strings"

	"github.com/Edditoria/grelease/github"
)

const (
	usageUrl      = "github.com/OWNER/REPO"
	usageFetchCmd = "grelease fetch --repo=" + usageUrl
)

func main() {
	fetch := flag.NewFlagSet("fetch", flag.ExitOnError)
	fetchRepo := fetch.String("repo", "", "URL to fetch: "+usageUrl)

	if len(os.Args) < 2 {
		os.Stderr.WriteString("Usage: " + usageFetchCmd + "\n")
		os.Exit(1)
	}

	switch os.Args[1] {
	case fetch.Name():
		fetch.Parse(os.Args[2:])
		exitCode := handleFetch(fetchRepo)
		os.Exit(exitCode)
	default:
		os.Stderr.WriteString("grelease: invalid command\nUsage: " + usageFetchCmd + "\n")
		os.Exit(1)
	}
}

func handleFetch(repoFlag *string) (exitCode int) {
	exitCode = 1

	splited := strings.Split(*repoFlag, "/")
	if len(splited) != 3 {
		os.Stderr.WriteString("grelease: invalid repository url\nExpect: " + usageUrl + "\n")
		return
	}
	if splited[0] != "github.com" {
		os.Stderr.WriteString("grelease: hostname not supported currently\nExpect: " + usageUrl + "\n")
		return
	}

	repo := github.Repo{Owner: splited[1], Name: splited[2]}
	// TODO: fetch all releases here:
	releases, _, err := repo.ListReleases(1)
	if err != nil {
		panic(err)
	}
	repo.Releases = releases
	err = repo.WriteJson(os.Stdout, "", "\t")
	if err != nil {
		panic(err)
	}

	return 0
}
