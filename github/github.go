package github

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
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

// Fetch from Github. **Append** results to current `repo.Releases`.
func (repo *Repo) FetchReleasesOnce(perPage, page int) error {
	rUrl := "https://api.github.com/repos/" + repo.Owner + "/" + repo.Name +
		"/releases?per_page=" + strconv.Itoa(perPage) + "&page=" + strconv.Itoa(page)
	resp, err := http.Get(rUrl)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status + " from " + rUrl)
	}
	defer resp.Body.Close()
	bodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var releases []*Release
	if err := json.Unmarshal(bodyByte, &releases); err != nil {
		return err
	}
	repo.Releases = append(repo.Releases, releases...)
	return nil
}
