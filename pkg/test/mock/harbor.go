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

func (r *MockRequester) Get(url string, auth models.Auth, printBody bool) ([]byte, error) {
	body := []byte(r.JsonBody)
	return body, nil
}

func (r *MockRequester) Delete(url string, auth models.Auth) (bool, error) {
	return false, nil
}

// mock harbor

type MockHarbor struct {
}

func (h *MockHarbor) GetAll() []models.HarborProject {
	projects := []models.HarborProject{}

	// get projects
	projectBody := []byte(payload.GetHarborProjects())
	_, projects_page := ParseProjects(projectBody, false)
	projects = append(projects, projects_page...)

	// iterate over projects
	for pIndex, project := range projects {
		if project.Name == "web" {
			// get repos
			repoBody := []byte(payload.GetHarborRepos())
			_, repos := ParseRepos(repoBody, false)

			if len(repos) > 0 {
				projects[pIndex].Repos = append(projects[pIndex].Repos, repos...)
			}
		}

		// get artifacts
		for rIndex, _ := range projects[pIndex].Repos {
			// repoName := GetRepoName(repo.Name)
			// if repoName == "app" {
			// fmt.Println("Get", project.Name, repoName)
			artifactBody := []byte(payload.GetHarborArtifacts())
			_, artifacts := ParseArtifacts(artifactBody, false)

			if len(artifacts) > 0 {
				projects[pIndex].Repos[rIndex].Artifacts = append(projects[pIndex].Repos[rIndex].Artifacts, artifacts...)
			}
			// }
		}
	}

	// for _, project := range projects {
	// 	fmt.Println(project.Name, "has repos", strconv.Itoa(len(project.Repos)))
	// }

	return projects
}

// helper - duplicated code!

func GetRepoName(fullName string) string {
	return strings.Split(fullName, "/")[1]
}

func ParseProjects(body []byte, print bool) (error, []models.HarborProject) {
	// parse body
	var items []models.HarborProject
	if err := json.Unmarshal(body, &items); err != nil {
		fmt.Println("Cannot unmarshal json")
		fmt.Println(string(body))
		fmt.Println(err)
		return err, nil
	}
	if print {
		fmt.Println(npjson.PrettyPrint(items))
	}

	return nil, items
}

func ParseRepos(body []byte, print bool) (error, []models.HarborRepo) {
	// parse body
	var items []models.HarborRepo
	if err := json.Unmarshal(body, &items); err != nil {
		fmt.Println("Cannot unmarshal json")
		fmt.Println(string(body))
		fmt.Println(err)
		return err, nil
	}
	if print {
		fmt.Println(npjson.PrettyPrint(items))
	}

	return nil, items
}

func ParseArtifacts(body []byte, print bool) (error, []models.HarborArtifact) {
	// parse body
	var items []models.HarborArtifact
	if err := json.Unmarshal(body, &items); err != nil {
		fmt.Println("Cannot unmarshal json")
		fmt.Println(string(body))
		fmt.Println(err)
		return err, nil
	}
	if print {
		fmt.Println(npjson.PrettyPrint(items))
	}

	return nil, items
}
