package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"

	"github.com/Edditoria/grelease/github"
)

const (
	usageUrl       = "github.com/OWNER/REPO"
	usageFetchCmd  = "grelease fetch [-p] " + usageUrl + " FILE_TO_WRITE"
	usageFetchCmd0 = "grelease fetch " + usageUrl + " FILE_TO_WRITE"
	usageFetchCmd1 = "grelease fetch -p " + usageUrl
)

func main() {
	fetch := flag.NewFlagSet("fetch", flag.ExitOnError)
	fetchPrintFlag := fetch.Bool("p", false, "print to stdout (will reject FILE_TO_WRITE)")

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
		exitCode := handleFetch(fetch, fetchPrintFlag)
		os.Exit(exitCode)
	default:
		os.Stderr.WriteString("grelease: invalid command\n")
		os.Stderr.WriteString("Usage: " + usageFetchCmd + "\n")
		os.Exit(1)
	}
}

// Subcommand "fetch" to consume API call, then write to a file or `stdout`.
// It also handle errors for users.
//
// @return int for exit code.
func handleFetch(fSet *flag.FlagSet, printFlag *bool) int {

	// Check structure of fetch command:

	args := fSet.Args()
	if len(args) < 1 {
		os.Stderr.WriteString("grelease: missing operand\n")
		os.Stderr.WriteString("Usage: " + usageFetchCmd0 + "\n")
		os.Stderr.WriteString("   or: " + usageFetchCmd1 + "\nOptions:\n")
		fSet.PrintDefaults()
		return 1
	}
	if *printFlag && len(args) > 1 {
		os.Stderr.WriteString("grelease: command conflict with -p option\n")
		os.Stderr.WriteString("Expect: " + usageFetchCmd1 + "\n")
		return 1
	}
	if !*printFlag && len(args) != 2 {
		os.Stderr.WriteString("grelease: invalid command\n")
		os.Stderr.WriteString("Usage: " + usageFetchCmd0 + "\n")
		os.Stderr.WriteString("   or: " + usageFetchCmd1 + "\nOptions:\n")
		fSet.PrintDefaults()
		return 1
	}
	if len(args) > 2 {
		os.Stderr.WriteString("grelease: too many arguments\nExpect: " + usageFetchCmd + "\n")
		return 1
	}
	argRepoUrl := args[0]
	argFilePath := args[1]

	// Check repo URL:

	splited := strings.Split(argRepoUrl, "/")
	if len(splited) != 3 {
		os.Stderr.WriteString("grelease: invalid repository url\n" +
			"- Expect: " + usageUrl + "\n" +
			"- Got:    " + argRepoUrl + "\n")
		return 1
	}
	if splited[0] != "github.com" {
		os.Stderr.WriteString("grelease: hostname not supported in current version\n" +
			"- Expect: " + usageUrl + "\n" +
			"- Got:    " + argRepoUrl + "\n")
		return 1
	}

	// Fetch releases:

	repo := github.Repo{Owner: splited[1], Name: splited[2]}
	// TODO: fetch all releases here:
	releases, _, err := repo.ListReleases(1)
	if err != nil {
		panic(err)
	}
	repo.Releases = releases

	// Write to stdout or file:

	if *printFlag {
		err = repo.WriteJson(os.Stdout, "", "\t")
		if err != nil {
			os.Stderr.WriteString(err.Error() + "\n")
			return 1
		}
		return 0
	}
	absPath, err := filepath.Abs(argFilePath)
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
