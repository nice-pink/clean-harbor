package mock

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nice-pink/clean-harbor/pkg/models"
	"github.com/nice-pink/clean-harbor/pkg/test/payload"
	npjson "github.com/nice-pink/goutil/pkg/json"
)

// requester

type MockRequester struct {
	JsonBody string
	Err      error
}

func (r *MockRequester) Get(url string, printBody bool) ([]byte, error) {
	body := []byte(r.JsonBody)
	return body, nil
}

func (r *MockRequester) Delete(url string) (bool, error) {
	return false, nil
}

// mock harbor

type MockHarbor struct {
	DeleteSuccessful bool
	DeleteError      error
}

func (h *MockHarbor) GetAll() map[string]models.RegistryProject {
	projects := map[string]models.RegistryProject{}

	// get projects
	projectBody := []byte(payload.GetHarborProjects())
	projects_page, _ := ParseProjects(projectBody, false)
	for _, project := range projects_page {
		project.Repos = map[string]models.HarborRepo{}
		projects[project.Name] = models.RegistryProject{
			Name:      project.Name,
			Id:        project.Id,
			RepoCount: project.RepoCount,
			Repos:     map[string]models.RegistryRepo{},
		}
	}

	// iterate over projects
	for pIndex, project := range projects {
		if project.Name == "web" {
			// get repos
			repoBody := []byte(payload.GetHarborRepos())
			harborRepos, _ := ParseRepos(repoBody, false)
			repos := GetRegistryReposFromArray(harborRepos)

			if len(repos) > 0 {
				// projects[project.Name].Repos = append(projects[project.Name].Repos, repos...)
				for _, repo := range repos {
					projects[project.Name].Repos[repo.Name] = repo
				}
			}
		}

		// get artifacts
		for _, repo := range projects[pIndex].Repos {
			artifactBody := []byte(payload.GetHarborArtifacts())
			harborArtifacts, _ := ParseArtifacts(artifactBody, false)
			artifacts := GetRegistryArtifacts(harborArtifacts)

			if len(artifacts) > 0 {
				repo.Artifacts = append(repo.Artifacts, artifacts...)
			}
			// }
			projects[project.Name].Repos[repo.Name] = repo
		}

	}

	// for _, project := range projects {
	// 	fmt.Println(project.Name, "has repos", strconv.Itoa(len(project.Repos)))
	// }

	return projects
}

func (h *MockHarbor) GetAllRepos(projectName string, print bool) (map[string]models.RegistryRepo, error) {
	// request
	repos := map[string]models.HarborRepo{}

	// iterate over projects
	if projectName == "web" {
		// get repos
		repoBody := []byte(payload.GetHarborRepos())
		pRepos, _ := ParseRepos(repoBody, false)
		for _, repo := range pRepos {
			repos[repo.Name] = repo
		}
	}

	// get artifacts
	for _, repo := range repos {
		artifactBody := []byte(payload.GetHarborArtifacts())
		artifacts, _ := ParseArtifacts(artifactBody, false)

		if len(artifacts) > 0 {
			repo.Artifacts = append(repo.Artifacts, artifacts...)
		}
		// }
		repos[repo.Name] = repo
	}

	// for _, project := range projects {
	// 	fmt.Println(project.Name, "has repos", strconv.Itoa(len(project.Repos)))
	// }

	return GetRegistryRepos(repos), nil
}

func (h *MockHarbor) EnrichReposWithArtificats(projects map[string]models.RegistryProject) map[string]models.RegistryProject {
	return projects
}

func (h *MockHarbor) DeleteArtifact(artifactReference string, projectName string, repoName string) (bool, error) {
	if h.DeleteSuccessful {
		return true, h.DeleteError
	}
	return true, h.DeleteError
}

func (h *MockHarbor) DeleteRepo(projectName string, repoName string) (bool, error) {
	if h.DeleteSuccessful {
		return true, h.DeleteError
	}
	return true, h.DeleteError
}

// helper - duplicated code!

func GetRepoName(fullName string) string {
	return strings.Split(fullName, "/")[1]
}

func ParseProjects(body []byte, print bool) ([]models.HarborProject, error) {
	// parse body
	var items []models.HarborProject
	if err := json.Unmarshal(body, &items); err != nil {
		fmt.Println("Cannot unmarshal json")
		fmt.Println(string(body))
		fmt.Println(err)
		return nil, err
	}
	if print {
		fmt.Println(npjson.PrettyPrint(items))
	}

	return items, nil
}

func ParseRepos(body []byte, print bool) ([]models.HarborRepo, error) {
	// parse body
	var items []models.HarborRepo
	if err := json.Unmarshal(body, &items); err != nil {
		fmt.Println("Cannot unmarshal json")
		fmt.Println(string(body))
		fmt.Println(err)
		return nil, err
	}
	if print {
		fmt.Println(npjson.PrettyPrint(items))
	}

	// fix repo name: from project/repo -> repo
	for index, item := range items {
		items[index].Name = GetRepoName(item.Name)
	}

	return items, nil
}

func ParseArtifacts(body []byte, print bool) ([]models.HarborArtifact, error) {
	// parse body
	var items []models.HarborArtifact
	if err := json.Unmarshal(body, &items); err != nil {
		fmt.Println("Cannot unmarshal json")
		fmt.Println(string(body))
		fmt.Println(err)
		return nil, err
	}
	if print {
		fmt.Println(npjson.PrettyPrint(items))
	}

	return items, nil
}

//

func GetRegistryRepos(harborRepos map[string]models.HarborRepo) map[string]models.RegistryRepo {
	repos := map[string]models.RegistryRepo{}

	for _, repo := range harborRepos {
		repos[repo.Name] = models.RegistryRepo{
			Name:          repo.Name,
			ArtifactCount: repo.ArtifactCount,
			Id:            repo.Id,
			Artifacts:     GetRegistryArtifacts(repo.Artifacts)}
	}

	return repos
}

func GetRegistryReposFromArray(harborRepos []models.HarborRepo) map[string]models.RegistryRepo {
	repos := map[string]models.RegistryRepo{}

	for _, repo := range harborRepos {
		repos[repo.Name] = models.RegistryRepo{
			Name:          repo.Name,
			ArtifactCount: repo.ArtifactCount,
			Id:            repo.Id,
			Artifacts:     GetRegistryArtifacts(repo.Artifacts)}
	}

	return repos
}

func GetRegistryArtifacts(harborArtifacts []models.HarborArtifact) []models.RegistryArtifact {
	artifacts := []models.RegistryArtifact{}
	for _, hArtifact := range harborArtifacts {
		tags := []models.RegistryTag{}
		for _, tag := range hArtifact.Tags {
			tags = append(tags, models.RegistryTag{Name: tag.Name, Created: tag.Created})
		}

		artifacts = append(artifacts, models.RegistryArtifact{
			Tags:    tags,
			ID:      hArtifact.ID,
			Digest:  hArtifact.Digest,
			Created: hArtifact.Created,
		})
	}
	return artifacts
}
