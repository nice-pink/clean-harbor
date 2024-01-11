package cleaner

import (
	"testing"

	"github.com/nice-pink/clean-harbor/pkg/test/mock"
	// "github.com/nice-pink/clean-harbor/pkg/test/payload"
)

// all

func TestFindUnsued(t *testing.T) {
	h := &mock.MockHarbor{}
	TAGS_HISTORY := 1

	c := NewCleaner(h, TAGS_HISTORY)
	extensions := []string{".yaml"}
	unused := c.FindUnused("../../pkg/test/repo", "repo.url", extensions)

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
	want_r2_name := "app-feature-1315"

	if got_r2_name != want_r2_name {
		t.Errorf("got_r2_name %q != want_r2_name %q", got_r2_name, want_r2_name)
	}

	got_t2 := unused[0].Projects[1].Repos[1].Tags
	want_t2 := []string{"0_0_zzzzzz"}

	index = 0
	for _, item := range got_t2 {
		if item != want_t2[index] {
			t.Errorf("got_t2 %q != want_t2 %q", got_t2, want_t2)
		}
		index++
	}
}

func TestFindUnsuedNone(t *testing.T) {
	h := &mock.MockHarbor{}
	TAGS_HISTORY := 5

	c := NewCleaner(h, TAGS_HISTORY)
	extensions := []string{".yaml"}
	unused := c.FindUnused("../../pkg/test/repo", "repo.url", extensions)

	got_base_name := unused[0].Name
	want_base_name := "repo.url"

	if got_base_name != want_base_name {
		t.Errorf("got_base_name %q != want_base_name %q", got_base_name, want_base_name)
	}

	// no artifacts
	got_t2 := unused[0].Projects[1].Repos[1].Tags
	want_t2 := []string{}

	index := 0
	for _, item := range got_t2 {
		if item != want_t2[index] {
			t.Errorf("got_t2 %q != want_t2 %q", got_t2, want_t2)
		}
		index++
	}
}
