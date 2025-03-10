package github

import (
	"testing"
)

func TestFetchReleaseOnce(t *testing.T) {
	hugo := Repo{Owner: "gohugoio", Name: "hugo"}
	perPage := 2
	err := hugo.FetchReleasesOnce(perPage, 1)
	if err != nil {
		t.Fatal(err)
	}
	rLen := len(hugo.Releases)
	if len(hugo.Releases) != rLen {
		t.Fatalf("len(hugo.Releases) expected %d but got %d", perPage, rLen)
	}
}
