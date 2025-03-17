package github

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Constants defined by Github API.
const (
	ApiVersion string = "2022-11-28"
	ApiTimeout time.Duration = 10 * time.Second
	HeaderAcceptGithubJson string = "application/vnd.github+json"
	HeaderApiVersionKey string = "X-GitHub-Api-Version"

	ReleasesPerPageDefault int = 30
	ReleasesPerPageMax     int = 100
)

// Get what we need only.
type Repo struct {
	Owner    string     `json:"owner"`
	Name     string     `json:"name"`
	Releases []*Release `json:"releases"`
}

// Get what we need only.
type Release struct {
	Id         int            `json:"id"`
	TagName    string         `json:"tag_name"`
	Name       string         `json:"name"`
	Draft      bool           `json:"draft"`
	Prerelease bool           `json:"prerelease"`
	Assets     []ReleaseAsset `json:"assets"`
}

// Get what we need only.
type ReleaseAsset struct {
	Id   int    `json:"id"`
	Url  string `json:"url"`
	Name string `json:"name"`
	Size int    `json:"size"`
}

func (repo *Repo) ReleasesUrl(perPage, page int) (*url.URL, error) {
	if perPage < 1 || perPage > ReleasesPerPageMax {
		msg := "perPage must be between 1 and " + strconv.Itoa(ReleasesPerPageMax)
		return nil, errors.New(msg)
	}
	if page < 1 {
		return nil, errors.New("page must be larger than 0")
	}
	rUrl := "https://api.github.com/repos/" + repo.Owner + "/" + repo.Name +
		"/releases?per_page=" + strconv.Itoa(perPage) + "&page=" + strconv.Itoa(page)
	return url.Parse(rUrl)
}

// List releases for a selected page. To fetch all releases, please do [Repo.ReloadReleases].
//
// It fetches 100 releases per page, instead of 30 that is default in Github API.
// The returned response pointer may be useful for [GetMaxPage] and debug.
func (repo *Repo) ListReleases(page int) ([]*Release, *http.Response, error) {
	rUrl, err := repo.ReleasesUrl(ReleasesPerPageMax, page)
	if err != nil {
		return nil, nil, err
	}
	rUrlStr := rUrl.String()
	req, err := http.NewRequest(http.MethodGet, rUrlStr, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("Accept", HeaderAcceptGithubJson)
	req.Header.Add(HeaderApiVersionKey, ApiVersion)
	client := http.Client{
		Timeout: ApiTimeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, nil, errors.New(resp.Status + " from " + rUrlStr)
	}
	defer resp.Body.Close()
	bodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	var releases []*Release
	if err := json.Unmarshal(bodyByte, &releases); err != nil {
		return nil, nil, err
	}
	return releases, resp, nil
}
