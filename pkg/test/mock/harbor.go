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
}

func (h *MockHarbor) GetAll() map[string]models.HarborProject {
	projects := map[string]models.HarborProject{}

	// get projects
	projectBody := []byte(payload.GetHarborProjects())
	projects_page, _ := ParseProjects(projectBody, false)
	for _, project := range projects_page {
		project.Repos = map[string]models.HarborRepo{}
		projects[project.Name] = project
	}

	// iterate over projects
	for pIndex, project := range projects {
		if project.Name == "web" {
			// get repos
			repoBody := []byte(payload.GetHarborRepos())
			repos, _ := ParseRepos(repoBody, false)

			if len(repos) > 0 {
				// projects[project.Name].Repos = append(projects[project.Name].Repos, repos...)
				for _, repo := range repos {
					projects[project.Name].Repos[repo.Name] = repo
				}
			}
		}

		// get artifacts
		for _, repo := range projects[pIndex].Repos {
			// repoName := GetRepoName(repo.Name)
			// if repoName == "app" {
			// fmt.Println("Get", project.Name, repoName)
			artifactBody := []byte(payload.GetHarborArtifacts())
			artifacts, _ := ParseArtifacts(artifactBody, false)

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

func (h *MockHarbor) GetAllRepos(projectName string, print bool) (map[string]models.HarborRepo, error) {
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
		// repoName := GetRepoName(repo.Name)
		// if repoName == "app" {
		// fmt.Println("Get", project.Name, repoName)
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

	return repos, nil
}

func EnrichReposWithArtificats(projects map[string]models.HarborProject) map[string]models.HarborProject {
	return projects
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
