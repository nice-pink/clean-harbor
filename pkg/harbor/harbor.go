package harbor

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/nice-pink/clean-harbor/pkg/models"
)

// dependencies

type Requester interface {
	Get(url string, auth models.Auth, printBody bool) ([]byte, error)
	Delete(url string, auth models.Auth) (bool, error)
}

// harbor

type HarborConfig struct {
	DryRun    bool
	HarborUrl string
	BasicAuth models.Auth
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

	// get projects
	index := 1
	for true {
		err, projects_page := h.GetProjects(index, 100, false)
		if err != nil {
			return nil
		}
		if len(projects_page) > 0 {
			projects = append(projects, projects_page...)
		} else {
			fmt.Println("Break", strconv.Itoa(len(projects)))
			break
		}
		index++
	}

	// iterate over projects
	for pIndex, project := range projects {
		// get repos
		index = 1
		for true {
			err, repos := h.GetRepos(project.Name, index, 100, false)
			if err != nil {
				continue
			}

			if len(repos) > 0 {
				projects[pIndex].Repos = append(projects[pIndex].Repos, repos...)
			} else {
				// fmt.Println("Got repos:", strconv.Itoa(len(projects[pIndex].Repos)))
				break
			}
			index++
		}

		// get artifacts
		for rIndex, repo := range projects[pIndex].Repos {
			index = 1
			for true {
				// fmt.Println("Get", project.Name, repo.Name)
				repoName := strings.Split(repo.Name, "/")[1]
				err, artifacts := h.GetArtifacts(project.Name, repoName, index, 100, false)
				if err != nil {
					continue
				}

				if len(artifacts) > 0 {
					projects[pIndex].Repos[rIndex].Artifacts = append(projects[pIndex].Repos[rIndex].Artifacts, artifacts...)
				} else {
					fmt.Println(repo.Name, " has artifacts: ", strconv.Itoa(len(projects[pIndex].Repos[rIndex].Artifacts)))
					break
				}
				index++
			}
		}
	}

	// for _, project := range projects {
	// 	fmt.Println(project.Name, "has repos", strconv.Itoa(len(project.Repos)))
	// }

	return projects
}

// project

func (h *Harbor) GetProjects(page int, pageSize int, print bool) (error, []models.HarborProject) {
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
	if print {
		fmt.Println(PrettyPrint(items))
	}

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

func (h *Harbor) GetRepos(projectName string, page int, pageSize int, print bool) (error, []models.HarborRepo) {
	// request
	path := "/projects/" + projectName + "/repositories" + h.GetQuery(page, pageSize)
	url := h.config.HarborUrl + path
	body, err := h.requester.Get(url, h.config.BasicAuth, false)
	if err != nil {
		fmt.Println("Could not request repo.")
		return err, nil
	}
	// fmt.Println(string(body))

	// parse body
	var items []models.HarborRepo
	if err := json.Unmarshal(body, &items); err != nil {
		fmt.Println("Cannot unmarshal json")
		fmt.Println(string(body))
		fmt.Println(err)
		return err, nil
	}
	if print {
		fmt.Println(PrettyPrint(items))
	}

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

func (h *Harbor) GetArtifacts(projectName string, repoName string, page int, pageSize int, print bool) (error, []models.HarborArtifact) {
	// request
	path := "/projects/" + projectName + "/repositories/" + repoName + "/artifacts" + h.GetQuery(page, pageSize)
	url := h.config.HarborUrl + path
	body, err := h.requester.Get(url, h.config.BasicAuth, false)
	if err != nil {
		fmt.Println("Could not request artifacts.")
		return err, nil
	}
	// fmt.Println(string(body))

	// parse body
	var items []models.HarborArtifact
	if err := json.Unmarshal(body, &items); err != nil {
		fmt.Println("Cannot unmarshal json")
		fmt.Println(string(body))
		fmt.Println(err)
		return err, nil
	}
	if print {
		fmt.Println(PrettyPrint(items))
	}

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

func BuildUniModels(projects map[string]models.HarborProject, baseUrl string) []models.UniBase {
	uBases := []models.UniBase{}
	uProjects := []models.UniProject{}
	// fmt.Println("Base:", base.Name)
	for _, project := range projects {
		uRepos := []models.UniRepo{}
		// fmt.Println("	Project:", project.Name)
		// fmt.Println("		", project.Name, "has repos", strconv.Itoa(len(project.Repos)))
		for _, repo := range project.Repos {
			// fmt.Println("		", repo.Name, repo.Tags)
			tags := []string{}
			for _, artifact := range repo.Artifacts {
				for _, tag := range artifact.Tags {
					tags = append(tags, tag.Name)
				}
			}
			uRepos = append(uRepos, models.UniRepo{Name: repo.Name, Tags: tags})
		}
		uProjects = append(uProjects, models.UniProject{Name: project.Name, Repos: uRepos})
	}
	uBases = append(uBases, models.UniBase{Name: baseUrl, Projects: uProjects})
	return uBases
}

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
