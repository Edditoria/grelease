package releasable

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/Edditoria/grelease/github"
)

func readRepoFromFile(repo *github.Repo, filePath string) error {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(absPath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &repo); err != nil {
		return err
	}
	return err
}

func TestInspect(t *testing.T) {
	var repoHugo100 github.Repo
	v333Re, err := regexp.Compile(Version333 + "$")
	if err != nil {
		t.Fatalf("pre-test: %v", err)
	}
	err = readRepoFromFile(&repoHugo100, "testdata/hugo_100_releases.json")
	if err != nil {
		t.Fatalf("pre-test: %v", err)
	}

	t.Run("hugo recent 100", func(t *testing.T) {
		opt := Option{
			IncludeDraft:      false,
			IncludePrerelease: false,
			TagNameRegex:      v333Re,
		}
		results, repoErr := Inspect(&repoHugo100, opt)
		if repoErr != nil {
			t.Fatalf("should not have error: %v", repoErr)
		}
		if len(results) != 100 {
			t.Fatalf("want 100 results but got %d", len(results))
		}
		var errResults []ResultForRelease
		for _, result := range results {
			if result.Error != nil {
				errResults = append(errResults, *result)
			}
		}
		wantLen := 1
		wantError := ErrReleaseNoAsset
		if len(errResults) == 0 {
			t.Errorf("release error: want %d records but got 0", wantLen)
		} else if len(errResults) != wantLen {
			t.Errorf("release error: want %d records but got %d\n"+
				"- 1st result: %s: %v\n"+
				"- 2nd result: %s: %v",
				wantLen, len(errResults),
				errResults[0].Release.TagName, errResults[0],
				errResults[1].Release.TagName, errResults[1])
		}
		if !errors.Is(errResults[0].Error, wantError) {
			t.Errorf("release error not match: (tag_name %s)\n"+
				"- want: %v\n"+
				"- got:  %v",
				errResults[0].Release.TagName, wantError, errResults[0].Error)
		}
	})
}
