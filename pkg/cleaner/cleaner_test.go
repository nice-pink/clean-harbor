package cleaner

import (
	"testing"

	"github.com/nice-pink/clean-harbor/pkg/test/mock"
	// "github.com/nice-pink/clean-harbor/pkg/test/payload"
)

// all

func TestFindUnsued(t *testing.T) {
	h := &mock.MockHarbor{}
	TAGS_HISTORY := 0
	UNKNOWN_HISTORY := 0

	c := NewCleaner(h, true, TAGS_HISTORY, UNKNOWN_HISTORY)
	extensions := []string{".yaml"}
	filterProjects := []string{}
	filterRepos := ""
	_, _, unused := c.FindUnused("../../pkg/test/repo", "repo.url", extensions, filterProjects, filterRepos, false, false)

	got_base_name := unused[0].Name
	want_base_name := "repo.url"

	if got_base_name != want_base_name {
		t.Errorf("got_base_name %q != want_base_name %q", got_base_name, want_base_name)
	}

	// 1 project
	got_p1_name := unused[0].Projects[0].Name
	want_p1_name := "dummy"

	if got_p1_name != want_p1_name {
		t.Errorf("got_p1_name %q != want_p1_name %q", got_p1_name, want_p1_name)
	}

	// 2 project
	got_p2_name := unused[0].Projects[1].Name
	want_p2_name := "web"

	if got_p2_name != want_p2_name {
		t.Errorf("got_p2_name %q != want_p1_name %q", got_p2_name, want_p2_name)
	}

	// 1 repo
	got_r1_name := unused[0].Projects[1].Repos[0].Name
	want_r1_name := "app-feature-1315"

	if got_r1_name != want_r1_name {
		t.Errorf("got_r1_name %q != want_r1_name %q", got_r1_name, want_r1_name)
	}

	got_t1 := unused[0].Projects[1].Repos[0].Tags
	want_t1 := []string{"0_0_xxxxx", "0_0_yyyyyy", "0_0_zzzzzz"}

	index := 0
	for _, item := range got_t1 {
		if item != want_t1[index] {
			t.Errorf("got_t1 %q != want_t1 %q", got_t1, want_t1)
		}
		index++
	}

	// 2 repo
	got_r2_name := unused[0].Projects[1].Repos[0].Name
	want_r2_name := "app"

	if got_r2_name != want_r2_name {
		t.Errorf("got_r2_name %q != want_r2_name %q", got_r2_name, want_r2_name)
	}

	got_t2 := unused[0].Projects[1].Repos[1].Tags
	want_t2 := []string{"0_0_zzzzzz"}

	index = 0
	for _, item := range want_t2 {
		if item != got_t2[index] {
			t.Errorf("%d got_t2 %q != want_t2 %q", index, got_t2, want_t2)
		}
		index++
	}
}

func TestFindUnsuedNone(t *testing.T) {
	// app is not contained in used!
	h := &mock.MockHarbor{}
	TAGS_HISTORY := 5
	UNKNOWN_HISTORY := 0

	c := NewCleaner(h, true, TAGS_HISTORY, UNKNOWN_HISTORY)
	extensions := []string{".yaml"}
	filterProjects := []string{}
	filterRepos := ""
	_, _, unused := c.FindUnused("../../pkg/test/repo", "repo.url", extensions, filterProjects, filterRepos, false, false)

	got_base_name := unused[0].Name
	want_base_name := "repo.url"

	if got_base_name != want_base_name {
		t.Errorf("got_base_name %q != want_base_name %q", got_base_name, want_base_name)
	}

	for _, project := range unused[0].Projects {
		if project.Name == "web" {
			for _, repo := range project.Repos {
				if repo.Name == "app" {
					t.Errorf("All images are used or within the tags_history, but found image unused. %s", repo.Name)
				}
			}
		}
	}
}

func TestFindUnsuedFilterProjects(t *testing.T) {
	h := &mock.MockHarbor{}
	TAGS_HISTORY := 5
	UNKNOWN_HISTORY := 0

	c := NewCleaner(h, true, TAGS_HISTORY, UNKNOWN_HISTORY)
	extensions := []string{}
	filterProjects := []string{"dummy"}
	filterRepos := ""
	_, _, unused := c.FindUnused("../../pkg/test/repo", "repo.url", extensions, filterProjects, filterRepos, false, false)

	got_base_name := unused[0].Name
	want_base_name := "repo.url"

	if got_base_name != want_base_name {
		t.Errorf("got_base_name %q != want_base_name %q", got_base_name, want_base_name)
	}

	// project cound
	got_len := len(unused[0].Projects)
	want_len := 1

	if got_len != want_len {
		t.Errorf("Filter project. Items in slice. got_len %q != want_len %q", got_len, want_len)
	}

	// no artifacts
	got_t2 := len(unused[0].Projects[0].Repos)
	want_t2 := 0

	if got_t2 != want_t2 {
		t.Errorf("Filter project. repos. got_t2 %q != want_t2 %q", got_t2, want_t2)
	}
}
