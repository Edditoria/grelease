package github

import (
	"net/http"
	"testing"
)

func TestReleasesUrl(t *testing.T) {
	want := `https://api.github.com/repos/gohugoio/hugo/releases?per_page=30&page=20`
	repo := Repo{Owner: "gohugoio", Name: "hugo"}
	rUrl, err := repo.ReleasesUrl(ReleasesPerPageDefault, 20)
	if err != nil {
		t.Fatal(err)
	}
	rUrlStr := rUrl.String()
	if rUrlStr != want {
		t.Errorf("repo.ReleasesUrl() not match:\n"+
			"- want: %s\n"+
			"- got:  %s",
			want, rUrlStr)
	}
}

func TestListReleases(t *testing.T) {
	hugo := Repo{Owner: "gohugoio", Name: "hugo"}
	releases, _, err := hugo.ListReleases(1)
	if err != nil {
		t.Fatal(err)
	}
	want := 100
	rLen := len(releases)
	if rLen != want {
		t.Fatalf("len(hugo.Releases) wants %d but got %d", want, rLen)
	}
}

func TestGetMaxPage(t *testing.T) {
	type testCase struct {
		name   string
		sample string
		want   int
	}
	tc0 := testCase{
		name:   "empty string",
		sample: "",
		want:   0,
	}
	tc1 := testCase{
		name:   "example from Github Docs",
		sample: `<https://api.github.com/repositories/1300192/issues?page=2>; rel="prev", <https://api.github.com/repositories/1300192/issues?page=4>; rel="next", <https://api.github.com/repositories/1300192/issues?page=515>; rel="last", <https://api.github.com/repositories/1300192/issues?page=1>; rel="first"`,
		want:   515,
	}
	tc2 := testCase{
		name:   "no link at all",
		sample: "[no link]",
		want:   0,
	}
	testCases := []testCase{tc0, tc1, tc2}
	for _, tc := range testCases {
		t.Run(tc.name, func(subT *testing.T) {
			header := make(http.Header)
			if tc.sample != "[no link]" {
				header.Set("link", tc.sample)
			}
			got := GetMaxPage(&header)
			if got != tc.want {
				subT.Errorf("wants %d but got %d", tc.want, got)
			}
		})
	}
}

func TestUpdateReleases(t *testing.T) {
	repo := Repo{Owner: "gohugoio", Name: "hugo"}
	err := repo.UpdateReleases()
	if err != nil {
		t.Fatal(err)
	}
	want := 316
	got := len(repo.Releases)
	if got < want {
		t.Fatalf("wants at least %d releases, but got %d", want, got)
	}
}
