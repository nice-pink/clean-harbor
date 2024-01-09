package harbor

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/nice-pink/clean-harbor/pkg/models"
)

var (
	config = HarborConfig{}
)

type HarborConfig struct {
	DryRun         bool
	HarborUrl      string
	HarborUser     string
	HarborPassword string
}

// config

func Configure(harborConfig HarborConfig) {
	config = harborConfig
}

// all

func GetAll() []models.HarborProject {
	err, projects := GetProjects(1, 100)
	if err != nil {
		fmt.Println(err)
		return nil
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

func GetProjects(page int, pageSize int) (error, []models.HarborProject) {
	// request
	path := "/projects" + GetQuery(page, pageSize)
	body, err := Get(path, false)
	if err != nil {
		fmt.Println("Could not request projects.")
		return err, nil
	}

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

func GetProject(id string) error {
	// request
	path := "/projects/" + id
	body, err := Get(path, false)
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

func GetRepos(projectName string, page int, pageSize int) (error, []models.HarborRepo) {
	// request
	path := "/projects/" + projectName + "/repositories" + GetQuery(page, pageSize)
	body, err := Get(path, false)
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

func GetRepo(name string, projectName string) error {
	path := "/projects/" + projectName + "/repositories/" + name
	body, err := Get(path, false)
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

func DeleteRepo(projectName string, repoName string) (bool, error) {
	path := "/projects/" + projectName + "/repositories/" + repoName
	if config.DryRun {
		fmt.Println("Delete:", path)
		return false, nil
	}
	success, err := Delete(path)
	if !success || err != nil {
		fmt.Println("Deleting not successful!")
	}
	return success, err
}

// artifact

func GetArtifacts(projectName string, repoName string, page int, pageSize int) (error, []models.HarborArtifact) {
	// request
	path := "/projects/" + projectName + "/repositories/" + repoName + "/artifacts" + GetQuery(page, pageSize)
	body, err := Get(path, false)
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

func GetArtifact(artifactReference string, projectName string, repoName string) error {
	path := "/projects/" + projectName + "/repositories/" + repoName + "/artifacts/" + artifactReference
	body, err := Get(path, false)
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

func DeleteArtifact(artifactReference string, projectName string, repoName string) (bool, error) {
	path := "/projects/" + projectName + "/repositories/" + repoName + "/artifacts/" + artifactReference
	success, err := Delete(path)
	if !success || err != nil {
		fmt.Println("Deleting not successful!")
	}
	return success, err
}

// request

func Get(path string, printBody bool) ([]byte, error) {
	url := config.HarborUrl + path
	// fmt.Println(url)

	// build request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// add basic auth
	req.SetBasicAuth(config.HarborUser, config.HarborPassword)

	// request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	// read and return
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	if printBody {
		fmt.Println(string(body))
	}

	return body, err
}

func Delete(path string) (bool, error) {
	url := config.HarborUrl + path
	// fmt.Println(url)

	// build request
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	// add basic auth
	req.SetBasicAuth(config.HarborUser, config.HarborPassword)

	// request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	defer resp.Body.Close()

	// read and return
	if resp.StatusCode != 200 && resp.StatusCode != 202 {
		fmt.Println("Could not delete. Status code:", strconv.Itoa(resp.StatusCode))
		return false, nil
	}
	return true, nil
}

// helper

func GetQuery(page int, pageSize int) string {
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
