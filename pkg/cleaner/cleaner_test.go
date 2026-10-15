package cleaner

import (
	"testing"

	"github.com/nice-pink/clean-harbor/pkg/models"
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
	_, _, unused, err := c.FindUnused("../../pkg/test/repo", "repo.url", extensions, filterProjects, filterRepos, false, false)
	if err != nil {
		t.Fatal(err)
	}

	got_base_name := unused[0].Name
	want_base_name := "repo.url"

	if got_base_name != want_base_name {
		t.Errorf("got_base_name %q != want_base_name %q", got_base_name, want_base_name)
	}

	// projects - look up by name, slice order is not deterministic (built from map)
	var webProject *models.UniProject
	foundDummy, foundWeb := false, false
	for i := range unused[0].Projects {
		switch unused[0].Projects[i].Name {
		case "dummy":
			foundDummy = true
		case "web":
			foundWeb = true
			webProject = &unused[0].Projects[i]
		}
	}
	if !foundDummy {
		t.Error("project 'dummy' not found in unused model")
	}
	if !foundWeb {
		t.Fatal("project 'web' not found in unused model")
	}

	// 1 repo - look up by name, slice order is not deterministic (built from map)
	var featureRepo *models.UniRepo
	for i := range webProject.Repos {
		if webProject.Repos[i].Name == "app-feature-1315" {
			featureRepo = &webProject.Repos[i]
			break
		}
	}
	if featureRepo == nil {
		t.Fatal("repo 'app-feature-1315' not found in unused model")
	}

	got_t1 := featureRepo.Tags
	want_t1 := []string{"0_0_xxxxx", "0_0_yyyyyy", "0_0_zzzzzz"}

	if len(got_t1) != len(want_t1) {
		t.Errorf("got_t1 %q != want_t1 %q", got_t1, want_t1)
	}
	for index, item := range want_t1 {
		if item != got_t1[index].Name {
			t.Errorf("%d got_t1 %q != want_t1 %q", index, got_t1, want_t1)
		}
	}

	// 2 repo - look up by name, slice order is not deterministic (built from map)
	var appRepo *models.UniRepo
	for i := range webProject.Repos {
		if webProject.Repos[i].Name == "app" {
			appRepo = &webProject.Repos[i]
			break
		}
	}
	if appRepo == nil {
		t.Fatal("repo 'app' not found in unused model")
	}

	got_t2 := appRepo.Tags
	want_t2 := []string{"0_0_zzzzzz"}

	if len(got_t2) != len(want_t2) {
		t.Errorf("got_t2 %q != want_t2 %q", got_t2, want_t2)
	}
	for index, item := range want_t2 {
		if item != got_t2[index].Name {
			t.Errorf("%d got_t2 %q != want_t2 %q", index, got_t2, want_t2)
		}
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
	_, _, unused, err := c.FindUnused("../../pkg/test/repo", "repo.url", extensions, filterProjects, filterRepos, false, false)
	if err != nil {
		t.Fatal(err)
	}

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
	_, _, unused, err := c.FindUnused("../../pkg/test/repo", "repo.url", extensions, filterProjects, filterRepos, false, false)
	if err != nil {
		t.Fatal(err)
	}

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

func TestFindUnusedUnknownRepoKeepsHistory(t *testing.T) {
	// an unknown repo (not referenced in manifests) with unknownHistory > 0
	// must NOT be queued for whole-repo deletion - only its older artifacts.
	h := &mock.MockHarbor{}
	TAGS_HISTORY := 0
	UNKNOWN_HISTORY := 1

	c := NewCleaner(h, true, TAGS_HISTORY, UNKNOWN_HISTORY)
	extensions := []string{".yaml"}
	filterProjects := []string{}
	filterRepos := ""
	unusedArtifacts, unusedRepos, unused, err := c.FindUnused("../../pkg/test/repo", "repo.url", extensions, filterProjects, filterRepos, false, false)
	if err != nil {
		t.Fatal(err)
	}

	// app-feature-1315 is unknown (not referenced in the test repo manifests)
	found := false
	for _, project := range unused[0].Projects {
		if project.Name != "web" {
			continue
		}
		for _, repo := range project.Repos {
			if repo.Name == "app-feature-1315" {
				found = true
				if repo.Unused {
					t.Error("unknown repo must not be marked Unused when unknownHistory > 0")
				}
				if len(repo.Tags) == 0 {
					t.Error("expected older tags of unknown repo to be queued for artifact deletion")
				}
			}
		}
	}
	if !found {
		t.Error("unknown repo app-feature-1315 not found in unused model")
	}

	for _, image := range unusedRepos {
		if image.Name == "app-feature-1315" {
			t.Error("unknown repo must not be queued for repo deletion when unknownHistory > 0")
		}
	}
	if len(unusedArtifacts) == 0 {
		t.Error("expected older artifacts of unknown repo to be queued for deletion")
	}
}

func TestFindUnusedUnknownRepoDeletedWithoutHistory(t *testing.T) {
	// with unknownHistory == 0 the whole unknown repo is queued for deletion
	h := &mock.MockHarbor{}
	TAGS_HISTORY := 0
	UNKNOWN_HISTORY := 0

	c := NewCleaner(h, true, TAGS_HISTORY, UNKNOWN_HISTORY)
	extensions := []string{".yaml"}
	filterProjects := []string{}
	filterRepos := ""
	_, unusedRepos, _, err := c.FindUnused("../../pkg/test/repo", "repo.url", extensions, filterProjects, filterRepos, false, false)
	if err != nil {
		t.Fatal(err)
	}

	found := false
	for _, image := range unusedRepos {
		if image.Name == "app-feature-1315" {
			found = true
		}
	}
	if !found {
		t.Error("unknown repo app-feature-1315 should be queued for repo deletion when unknownHistory == 0")
	}
}
