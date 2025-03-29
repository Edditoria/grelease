/*
Tool to inspect release data, that is fetched from Github API.
*/
package releasable

import (
	"errors"
	"regexp"

	"github.com/Edditoria/grelease/github"
)

const (
	Version333 = `^v[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}`
)

// Predefined errors to let you do [errors.Is].
var (
	ErrRepoNoRelease       = errors.New("repo has not release")
	ErrReleaseIsDraft      = errors.New("release is draft")
	ErrReleaseIsPrerelease = errors.New("release is pre-release")
	ErrReleaseNoAsset      = errors.New("release has not asset")
	ErrReleaseBadTagName   = errors.New("release has bad tag name")
)

type Option struct {
	IncludeDraft      bool // Should be false unless you know what you doing.
	IncludePrerelease bool // Should be false unless you want to include prereleases.
	TagNameRegex      *regexp.Regexp
}

type ResultForRelease struct {
	Release *github.Release
	Error   error
}

// TODO: Inspect assets in each release.
//   - @return One-on-one result for each release in a repo.
//   - @return Error that the repo has not any release.
func Inspect(repo *github.Repo, opt Option) ([]*ResultForRelease, error) {
	if len(repo.Releases) < 1 {
		return nil, ErrRepoNoRelease
	}
	var results []*ResultForRelease
	for _, rel := range repo.Releases {
		err := InspectRelease(rel, opt)
		results = append(results, &ResultForRelease{
			Release: rel,
			Error:   err,
		})
	}
	return results, nil
}

// Usage: If `err == nil`, that release is healthy.
// NOTE: It does not inspect assets in the release.
//
// @return Collection of error(s). Used [errors.Join].
func InspectRelease(release *github.Release, opt Option) error {
	var err1, err2, err3, err4 error
	if !opt.IncludeDraft && release.Draft {
		err1 = ErrReleaseIsDraft
	}
	if !opt.IncludePrerelease && release.Prerelease {
		err2 = ErrReleaseIsPrerelease
	}
	if len(release.Assets) < 1 {
		err3 = ErrReleaseNoAsset
	}
	if !opt.TagNameRegex.MatchString(release.TagName) {
		err4 = ErrReleaseBadTagName
	}
	return errors.Join(err1, err2, err3, err4)
}
