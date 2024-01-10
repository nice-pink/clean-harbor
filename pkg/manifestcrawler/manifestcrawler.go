package manifestcrawler

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	models "github.com/nice-pink/clean-harbor/pkg/models"
	"github.com/nice-pink/goutil/pkg/filesystem"
)

type Image struct {
	BaseUrl string
	Project string
	Name    string
	Tag     string
}

func GetImagesByRepo(folder string, repoUrl string, extensions []string) ([]string, error) {
	pattern := ".*(" + repoUrl + ".*)"
	replacement := "${1}"
	values, err := GetImages(folder, pattern, replacement, extensions, true)
	if err != nil {
		fmt.Println(err)
	}
	return values, err
}

func GetImages(folder string, pattern string, replacement string, extensions []string, recursive bool) ([]string, error) {
	values, err := filesystem.GetRegexInAllFiles(folder, true, pattern, replacement, extensions)
	if err != nil {
		fmt.Println(err)
	}

	images := []Image{}
	for _, value := range values {
		image := GetImage(value)
		// fmt.Println(image.Name)
		images = append(images, image)
		// break
	}

	GetImageProjects(images)

	return values, err
}

func GetImage(image string) Image {
	// fmt.Println(image)
	tag := ""
	if strings.Contains(image, ":") {
		// only get tag if there is one
		tag = strings.TrimSpace(filesystem.ReplaceRegex(image, "(.*):([a-zA-Z0-9_-]+)", "${2}"))
	}

	// fmt.Println(tag)
	container := strings.TrimSpace(filesystem.ReplaceRegex(image, "(.*):([a-zA-Z0-9_-]+)", "${1}"))
	name := strings.TrimSpace(filesystem.ReplaceRegex(container, ".*(/)([a-zA-Z0-9_-]+)", "${2}"))
	// fmt.Println(name)
	base := strings.TrimSpace(filesystem.ReplaceRegex(container, "(.*)(/)([a-zA-Z0-9_-]+)", "${1}"))
	project := strings.TrimSpace(filesystem.ReplaceRegex(base, "(.*)(/)([a-zA-Z0-9_-]+)", "${3}"))
	// fmt.Println(name)
	baseUrl := strings.TrimSpace(filesystem.ReplaceRegex(base, "(.*)(/)([a-zA-Z0-9_-]+)", "${1}"))
	// fmt.Println(name)
	return Image{Tag: tag, Name: name, Project: project, BaseUrl: baseUrl}
}

func GetImageProjects(images []Image) map[string]models.ManifestBase {
	bases := make(map[string]models.ManifestBase)

	// iterate over images
	for _, image := range images {
		//add base
		if base, ok := bases[image.BaseUrl]; !ok {
			// fmt.Println("Does not exist '"+image.BaseUrl+"'", image.Name, image.Project, image.Tag)
			// add project
			newProject := models.ManifestProject{Name: image.Project, Repos: make(map[string]models.ManifestRepo)}
			newProject.Repos[image.Name] = models.ManifestRepo{Name: image.Name, Tags: []string{image.Tag}}

			// create base
			newBase := models.ManifestBase{Name: image.BaseUrl, Projects: make(map[string]models.ManifestProject)}
			newBase.Projects[image.Project] = newProject

			// add base
			bases[image.BaseUrl] = newBase
		} else {
			// project
			if project, ok := base.Projects[image.Project]; !ok {
				// add project
				newProject := models.ManifestProject{Name: image.Project, Repos: make(map[string]models.ManifestRepo)}
				newProject.Repos[image.Name] = models.ManifestRepo{Name: image.Name, Tags: []string{image.Tag}}

				// create base
				base.Projects[image.Project] = newProject
			} else {
				// repo
				if repo, ok := project.Repos[image.Name]; !ok {
					// add project
					newRepo := models.ManifestRepo{Name: image.Name, Tags: []string{image.Tag}}
					project.Repos[image.Name] = newRepo
				} else {
					found := false
					for _, tag := range repo.Tags {
						if tag == image.Tag {
							found = true
							break
						}
					}
					// add tag
					if !found {
						repo.Tags = append(repo.Tags, image.Tag)
						project.Repos[image.Name] = repo
					}
				}
			}
		}
	}

	// fmt.Println("bases", strconv.Itoa(len(bases)))

	// for _, base := range bases {
	// 	fmt.Println("Base:", base.Name)
	// 	for _, project := range base.Projects {
	// 		fmt.Println("	Project:", project.Name)
	// 		// fmt.Println("		", project.Name, "has repos", strconv.Itoa(len(project.Repos)))
	// 		for _, repo := range project.Repos {
	// 			fmt.Println("		", repo.Name, repo.Tags)
	// 		}
	// 	}
	// }

	DumpJson(bases, "bin/repo.json")

	uBases := BuildUniModels(bases)
	DumpJson(uBases, "bin/easy.json")

	return bases
}

func BuildUniModels(bases map[string]models.ManifestBase) []models.UniBase {
	uBases := []models.UniBase{}
	for _, base := range bases {
		uProjects := []models.UniProject{}
		// fmt.Println("Base:", base.Name)
		for _, project := range base.Projects {
			uRepos := []models.UniRepo{}
			// fmt.Println("	Project:", project.Name)
			// fmt.Println("		", project.Name, "has repos", strconv.Itoa(len(project.Repos)))
			for _, repo := range project.Repos {
				// fmt.Println("		", repo.Name, repo.Tags)
				uRepos = append(uRepos, models.UniRepo{Name: repo.Name, Tags: repo.Tags})
			}
			uProjects = append(uProjects, models.UniProject{Name: project.Name, Repos: uRepos})
		}
		uBases = append(uBases, models.UniBase{Name: base.Name, Projects: uProjects})
	}
	return uBases
}

func DumpJson(i interface{}, filepath string) {

	j, _ := json.MarshalIndent(i, "", "  ")
	// fmt.Println(string(j))

	file, err := os.Create(filepath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	if _, err := file.Write(j); err != nil {
		fmt.Println(err)
	}
}