package github

import (
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

func TestFetchReleasesOnce(t *testing.T) {
	hugo := Repo{Owner: "gohugoio", Name: "hugo"}
	err := hugo.FetchReleasesOnce(1)
	if err != nil {
		t.Fatal(err)
	}
	want := 100
	rLen := len(hugo.Releases)
	if rLen != want {
		t.Fatalf("len(hugo.Releases) wants %d but got %d", want, rLen)
	}
}
