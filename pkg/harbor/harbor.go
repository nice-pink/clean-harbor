package harbor

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/nice-pink/clean-harbor/pkg/models"
	npjson "github.com/nice-pink/goutil/pkg/json"
	"github.com/nice-pink/goutil/pkg/log"
	"github.com/nice-pink/goutil/pkg/network"
)

// dependencies

type Requester interface {
	Get(url string, printBody bool) ([]byte, error)
	Delete(url string) (bool, error)
}

// harbor

type HarborConfig struct {
	DryRun    bool
	HarborUrl string
	BasicAuth network.Auth
}

type Harbor struct {
	requester Requester
	config    HarborConfig
}

func NewHarbor(requester Requester, config HarborConfig) *Harbor {
	return &Harbor{
		requester: requester,
		config:    config,
	}
}

// func (h *Harbor) Configure(harborConfig HarborConfig) {
// 	h.config = harborConfig
// }

// all

func (h *Harbor) GetAll(filterRepos string) map[string]models.RegistryProject {
	projects := map[string]models.RegistryProject{}

	// get projects
	index := 1
	for true {
		projects_page, err := h.GetProjects(index, 100, false)
		if err != nil {
			return nil
		}
		if len(projects_page) > 0 {
			for _, project := range projects_page {
				projects[project.Name] = models.RegistryProject{
					Name:      project.Name,
					Id:        project.Id,
					RepoCount: project.RepoCount,
					Repos:     map[string]models.RegistryRepo{},
				}
			}
		} else {
			fmt.Println(strconv.Itoa(len(projects)), "projects.")
			break
		}
		index++
	}

	// iterate over projects
	for _, project := range projects {
		// get repos
		index = 1
		for true {
			harborRepos, err := h.GetRepos(project.Name, index, 100, false)
			if err != nil {
				continue
			}
			repos := GetRegistryReposFromArray(harborRepos, filterRepos)

			if len(repos) > 0 {
				// projects[pIndex].Repos = append(projects[pIndex].Repos, repos...)
				// projects[pIndex].Repos = append(projects[pIndex].Repos, repos...)
				for _, repo := range repos {
					projects[project.Name].Repos[repo.Name] = repo
				}

			} else {
				fmt.Println(project.Name, "- repos:", strconv.Itoa(len(projects[project.Name].Repos)))
				break
			}
			index++
		}

		// get artifacts
		for _, repo := range projects[project.Name].Repos {
			index = 1
			for true {
				// fmt.Println("Get", project.Name, repo.Name)
				harborArtifacts, err := h.GetArtifacts(project.Name, repo.Name, index, 100, false)
				if err != nil {
					continue
				}
				artifacts := GetRegistryArtifacts(harborArtifacts)

				if len(artifacts) > 0 {
					repo.Artifacts = append(repo.Artifacts, artifacts...)
				} else {
					fmt.Println(repo.Name, " has artifacts: ", strconv.Itoa(len(repo.Artifacts)))
					break
				}
				index++
			}
			projects[project.Name].Repos[repo.Name] = repo
		}
	}

	// for _, project := range projects {
	// 	fmt.Println(project.Name, "has repos", strconv.Itoa(len(project.Repos)))
	// }

	return projects
}

// all

func (h *Harbor) EnrichReposWithArtificats(projects map[string]models.RegistryProject) map[string]models.RegistryProject {
	index := 1
	// iterate over projects
	for _, project := range projects {
		// get artifacts
		for _, repo := range projects[project.Name].Repos {
			index = 1
			for true {
				// fmt.Println("Get", project.Name, repo.Name)
				harborArtifacts, err := h.GetArtifacts(project.Name, repo.Name, index, 100, false)
				if err != nil {
					continue
				}
				artifacts := GetRegistryArtifacts(harborArtifacts)

				if len(artifacts) > 0 {
					repo.Artifacts = append(repo.Artifacts, artifacts...)
				} else {
					fmt.Println(repo.Name, " has artifacts: ", strconv.Itoa(len(repo.Artifacts)))
					break
				}
				index++
			}
			projects[project.Name].Repos[repo.Name] = repo
		}
	}

	// for _, project := range projects {
	// 	fmt.Println(project.Name, "has repos", strconv.Itoa(len(project.Repos)))
	// }

	return projects
}

func GetRepoName(fullName string) string {
	return strings.Split(fullName, "/")[1]
}

// project

func (h *Harbor) GetProjects(page int, pageSize int, print bool) ([]models.HarborProject, error) {
	// request
	path := "/projects" + h.GetQuery(page, pageSize)
	url := h.config.HarborUrl + path
	body, err := h.requester.Get(url, false)
	if err != nil {
		fmt.Println("Could not request projects.")
		return nil, err
	}
	// fmt.Println(string(body))

	return ParseProjects(body, print)
}

func GetRegistryProjects(harborProjects []models.HarborProject) []models.RegistryProject {
	projects := []models.RegistryProject{}

	for _, project := range harborProjects {
		projects = append(projects, models.RegistryProject{
			Name:      project.Name,
			Id:        project.Id,
			RepoCount: project.RepoCount,
			Repos:     GetRegistryRepos(project.Repos, ""),
		})
	}

	return projects
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

func (h *Harbor) GetProject(id string) error {
	// request
	path := "/projects/" + id
	url := h.config.HarborUrl + path
	body, err := h.requester.Get(url, false)
	if err != nil {
		fmt.Println("Could not request repo.")
		return err
	}

	// parse body
	return ParseProject(body)
}

func ParseProject(body []byte) error {
	// parse body
	var item models.HarborProject
	if err := json.Unmarshal(body, &item); err != nil {
		fmt.Println(string(body))
		fmt.Println("Cannot unmarshal json")
		return err
	}
	fmt.Println(npjson.PrettyPrint(item))

	return nil
}

// repo

func (h *Harbor) GetAllRepos(projectName string, filterRepos string, print bool) (map[string]models.RegistryRepo, error) {
	// request
	repos := map[string]models.HarborRepo{}

	index := 1
	for true {
		path := "/projects/" + projectName + "/repositories" + h.GetQuery(index, 100)
		url := h.config.HarborUrl + path

		body, err := h.requester.Get(url, false)
		if err != nil {
			fmt.Println("Could not request repo.")
			return nil, err
		}
		// fmt.Println(string(body))

		// parse body
		newRepos, err := ParseRepos(body, print)
		if len(newRepos) == 0 || err != nil {
			break
		}
		for _, repo := range newRepos {
			repos[repo.Name] = repo
		}
		index++
	}
	return GetRegistryRepos(repos, filterRepos), nil
}

func (h *Harbor) GetRepos(projectName string, page int, pageSize int, print bool) ([]models.HarborRepo, error) {
	// request
	path := "/projects/" + projectName + "/repositories" + h.GetQuery(page, pageSize)
	url := h.config.HarborUrl + path
	body, err := h.requester.Get(url, false)
	if err != nil {
		fmt.Println("Could not request repo.")
		return nil, err
	}
	// fmt.Println(string(body))

	// parse body
	return ParseRepos(body, print)
}

func GetRegistryRepos(harborRepos map[string]models.HarborRepo, filterRepos string) map[string]models.RegistryRepo {
	repos := map[string]models.RegistryRepo{}

	for _, repo := range harborRepos {
		if filterRepos == "" || strings.Contains(repo.Name, filterRepos) {
			repos[repo.Name] = models.RegistryRepo{
				Name:          repo.Name,
				ArtifactCount: repo.ArtifactCount,
				Id:            repo.Id,
				Artifacts:     GetRegistryArtifacts(repo.Artifacts)}
		}
	}

	return repos
}

func GetRegistryReposFromArray(harborRepos []models.HarborRepo, filterRepos string) map[string]models.RegistryRepo {
	repos := map[string]models.RegistryRepo{}

	for _, repo := range harborRepos {
		if filterRepos == "" || strings.Contains(repo.Name, filterRepos) {
			repos[repo.Name] = models.RegistryRepo{
				Name:          repo.Name,
				ArtifactCount: repo.ArtifactCount,
				Id:            repo.Id,
				Artifacts:     GetRegistryArtifacts(repo.Artifacts)}
		}
	}

	return repos
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

func (h *Harbor) GetRepo(name string, projectName string) error {
	path := "/projects/" + projectName + "/repositories/" + name
	url := h.config.HarborUrl + path
	body, err := h.requester.Get(url, false)
	if err != nil {
		fmt.Println("Could not request repo.")
		return err
	}

	// unmarshal body
	return ParseRepo(body)
}

func ParseRepo(body []byte) error {
	// unmarshal body
	var item models.HarborRepo
	if err := json.Unmarshal(body, &item); err != nil {
		fmt.Println("Cannot unmarshal json")
		fmt.Println(string(body))
		return err
	}
	fmt.Println(npjson.PrettyPrint(item))

	return nil
}

func (h *Harbor) DeleteRepo(projectName string, repoName string) (bool, error) {
	path := "/projects/" + projectName + "/repositories/" + repoName
	if h.config.DryRun {
		log.Info("Dry run: Delete:", path)
		return false, nil
	}
	url := h.config.HarborUrl + path
	success, err := h.requester.Delete(url)
	if !success || err != nil {
		log.Error("Deleting not successful!")
	}
	return success, err
}

// artifact

func (h *Harbor) GetArtifacts(projectName string, repoName string, page int, pageSize int, print bool) ([]models.HarborArtifact, error) {
	// request
	path := "/projects/" + projectName + "/repositories/" + repoName + "/artifacts" + h.GetQuery(page, pageSize)
	url := h.config.HarborUrl + path
	body, err := h.requester.Get(url, false)
	if err != nil {
		fmt.Println("Could not request artifacts.")
		return nil, err
	}
	// fmt.Println(string(body))

	// parse body
	return ParseArtifacts(body, print)
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

func (h *Harbor) GetArtifact(artifactReference string, projectName string, repoName string) error {
	path := "/projects/" + projectName + "/repositories/" + repoName + "/artifacts/" + artifactReference
	url := h.config.HarborUrl + path
	body, err := h.requester.Get(url, false)
	if err != nil {
		fmt.Println("Could not request repo.")
		return err
	}

	// unmarshal body
	return h.ParseArtifact(body)
}

func (h *Harbor) ParseArtifact(body []byte) error {
	// unmarshal body
	var item models.HarborArtifact
	if err := json.Unmarshal(body, &item); err != nil {
		fmt.Println("Cannot unmarshal json")
		fmt.Println(string(body))
		return err
	}
	fmt.Println(npjson.PrettyPrint(item))

	return nil
}

func (h *Harbor) DeleteArtifact(artifactReference string, projectName string, repoName string) (bool, error) {
	path := "/projects/" + projectName + "/repositories/" + repoName + "/artifacts/" + artifactReference
	url := h.config.HarborUrl + path

	if h.config.DryRun {
		log.Info("Dry run: Delete:", url)
		return false, nil
	}

	success, err := h.requester.Delete(url)
	if !success || err != nil {
		log.Error("Deleting not successful!")
	}
	return success, err
}

// helper

func (h *Harbor) GetQuery(page int, pageSize int) string {
	query := ""
	if page > 0 {
		query = "?page=" + strconv.Itoa(page)
	}
	if pageSize > 0 {
		if len(query) > 0 {
			query += "&"
		} else {
			query += "?"
		}
		query += "page_size=" + strconv.Itoa(pageSize)
	}
	return query
}
