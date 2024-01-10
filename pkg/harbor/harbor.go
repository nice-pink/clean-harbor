package harbor

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/nice-pink/clean-harbor/pkg/models"
	"github.com/nice-pink/clean-harbor/pkg/network"
)

// dependencies

type Requester interface {
	Get(url string, auth network.Auth, printBody bool) ([]byte, error)
	Delete(url string, auth network.Auth) (bool, error)
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

func (h *Harbor) GetAll() []models.HarborProject {
	projects := []models.HarborProject{}

	index := 1
	for true {
		err, projects_page := h.GetProjects(index, 100)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if len(projects_page) > 0 {
			projects = append(projects, projects_page...)
		} else {
			break
		}
		index++
	}

	// for pIndex, project := range projects {
	// 	err, repos := GetRepos(project.Name, 1, 100)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return nil
	// 	}

	// 	projects[pIndex].Repos = repos

	// 	// for _, repo := range repos {

	// 	// }
	// }

	// for _, project := range projects {
	// 	fmt.Print(project.Name, "has repos", strconv.Itoa(len(project.Repos)))
	// }

	return projects
}

// project

func (h *Harbor) GetProjects(page int, pageSize int) (error, []models.HarborProject) {
	// request
	path := "/projects" + h.GetQuery(page, pageSize)
	url := h.config.HarborUrl + path
	body, err := h.requester.Get(url, h.config.BasicAuth, false)
	if err != nil {
		fmt.Println("Could not request projects.")
		return err, nil
	}
	// fmt.Println(string(body))

	// parse body
	var items []models.HarborProject
	if err := json.Unmarshal(body, &items); err != nil {
		fmt.Println("Cannot unmarshal json")
		fmt.Println(string(body))
		fmt.Println(err)
		return err, nil
	}
	fmt.Println(PrettyPrint(items))

	return nil, items
}

func (h *Harbor) GetProject(id string) error {
	// request
	path := "/projects/" + id
	url := h.config.HarborUrl + path
	body, err := h.requester.Get(url, h.config.BasicAuth, false)
	if err != nil {
		fmt.Println("Could not request repo.")
		return err
	}

	// parse body
	var item models.HarborProject
	if err := json.Unmarshal(body, &item); err != nil {
		fmt.Println(string(body))
		fmt.Println("Cannot unmarshal json")
		return err
	}
	fmt.Println(PrettyPrint(item))

	return nil
}

// repo

func (h *Harbor) GetRepos(projectName string, page int, pageSize int) (error, []models.HarborRepo) {
	// request
	path := "/projects/" + projectName + "/repositories" + h.GetQuery(page, pageSize)
	url := h.config.HarborUrl + path
	body, err := h.requester.Get(url, h.config.BasicAuth, false)
	if err != nil {
		fmt.Println("Could not request repo.")
		return err, nil
	}

	// parse body
	var items []models.HarborRepo
	if err := json.Unmarshal(body, &items); err != nil {
		fmt.Println("Cannot unmarshal json")
		fmt.Println(string(body))
		fmt.Println(err)
		return err, nil
	}
	fmt.Println(PrettyPrint(items))

	return nil, items
}

func (h *Harbor) GetRepo(name string, projectName string) error {
	path := "/projects/" + projectName + "/repositories/" + name
	url := h.config.HarborUrl + path
	body, err := h.requester.Get(url, h.config.BasicAuth, false)
	if err != nil {
		fmt.Println("Could not request repo.")
		return err
	}

	// unmarshal body
	var item models.HarborRepo
	if err := json.Unmarshal(body, &item); err != nil {
		fmt.Println("Cannot unmarshal json")
		fmt.Println(string(body))
		return err
	}
	fmt.Println(PrettyPrint(item))

	return nil
}

func (h *Harbor) DeleteRepo(projectName string, repoName string) (bool, error) {
	path := "/projects/" + projectName + "/repositories/" + repoName
	if h.config.DryRun {
		fmt.Println("Delete:", path)
		return false, nil
	}
	url := h.config.HarborUrl + path
	success, err := h.requester.Delete(url, h.config.BasicAuth)
	if !success || err != nil {
		fmt.Println("Deleting not successful!")
	}
	return success, err
}

// artifact

func (h *Harbor) GetArtifacts(projectName string, repoName string, page int, pageSize int) (error, []models.HarborArtifact) {
	// request
	path := "/projects/" + projectName + "/repositories/" + repoName + "/artifacts" + h.GetQuery(page, pageSize)
	url := h.config.HarborUrl + path
	body, err := h.requester.Get(url, h.config.BasicAuth, false)
	if err != nil {
		fmt.Println("Could not request artifacts.")
		return err, nil
	}

	// parse body
	var items []models.HarborArtifact
	if err := json.Unmarshal(body, &items); err != nil {
		fmt.Println("Cannot unmarshal json")
		fmt.Println(string(body))
		fmt.Println(err)
		return err, nil
	}
	fmt.Println(PrettyPrint(items))

	return nil, items
}

func (h *Harbor) GetArtifact(artifactReference string, projectName string, repoName string) error {
	path := "/projects/" + projectName + "/repositories/" + repoName + "/artifacts/" + artifactReference
	url := h.config.HarborUrl + path
	body, err := h.requester.Get(url, h.config.BasicAuth, false)
	if err != nil {
		fmt.Println("Could not request repo.")
		return err
	}

	// unmarshal body
	var item models.HarborArtifact
	if err := json.Unmarshal(body, &item); err != nil {
		fmt.Println("Cannot unmarshal json")
		fmt.Println(string(body))
		return err
	}
	fmt.Println(PrettyPrint(item))

	return nil
}

func (h *Harbor) DeleteArtifact(artifactReference string, projectName string, repoName string) (bool, error) {
	path := "/projects/" + projectName + "/repositories/" + repoName + "/artifacts/" + artifactReference
	url := h.config.HarborUrl + path
	success, err := h.requester.Delete(url, h.config.BasicAuth)
	if !success || err != nil {
		fmt.Println("Deleting not successful!")
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

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
