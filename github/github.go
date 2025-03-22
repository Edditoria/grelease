package github

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// Constants defined by Github API.
const (
	ApiVersion             string        = "2022-11-28"
	ApiTimeout             time.Duration = 10 * time.Second
	HeaderAcceptGithubJson string        = "application/vnd.github+json"
	HeaderApiVersionKey    string        = "X-GitHub-Api-Version"

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

// Fetch all releases from Github. If all success, replace the old [Repo].Releases.
//
// @param maxCall to limit API call (page), while up to 100 releases per page.
func (repo *Repo) UpdateReleases(maxCall int) error {
	releases, resp, err := repo.ListReleases(1)
	if err != nil {
		return err
	}
	maxPage := GetMaxPage(&resp.Header)
	if maxPage < 2 {
		repo.Releases = releases
		return nil
	}
	if maxCall < maxPage {
		maxPage = maxCall
	}
	for i := 2; i <= maxPage; i++ {
		newReleases, _, err := repo.ListReleases(i)
		if err != nil {
			return err
		}
		releases = append(releases, newReleases...)
	}
	repo.Releases = releases
	return nil
}

// List releases for a selected page. To fetch all releases, please do [Repo.UpdateReleases].
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

// Get information of maximum page from Github response.
// If this function cannot find it, it returns 0.
//
// Basically this function parse the "link" field in the HTTP header.
// And search for URL with `rel="last"`.
//
// See: https://docs.github.com/en/rest/using-the-rest-api/using-pagination-in-the-rest-api
func GetMaxPage(header *http.Header) int {
	linkStr := header.Get("link")
	if linkStr == "" {
		// fmt.Println("GetMaxPage(): linkStr is empty")
		return 0
	}
	// We want: `<https://...?page=[X]>; rel="last"`
	parts := strings.SplitN(linkStr, `rel="last`, 2) // without tailing `"`.
	startIdx := 5 + strings.LastIndex(parts[0], "page=")
	endIdx := strings.LastIndex(parts[0], ">")
	if startIdx < 0 || endIdx < 0 {
		return 0
	}
	ans, err := strconv.Atoi(parts[0][startIdx:endIdx])
	if err != nil {
		return 0
	}
	return ans
}

func (repo *Repo) WriteJson(w io.Writer, prefix, indent string) error {
	enc := json.NewEncoder(w)
	enc.SetIndent(prefix, indent)
	err := enc.Encode(repo)
	return err
}

func (repo *Repo) WriteJsonFile(filePath, prefix, indent string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0o664)
	if err != nil {
		return err
	}
	defer file.Close()
	err = repo.WriteJson(file, prefix, indent)
	return err
}
