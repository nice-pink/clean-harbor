package manifestcrawler

import (
	"fmt"
	"os"
	"strings"

	"github.com/nice-pink/clean-harbor/pkg/models"
	"github.com/nice-pink/goutil/pkg/filesystem"
	"github.com/nice-pink/goutil/pkg/git"
)

var (
	ReposBaseFolder string = "bin/repo"
)

func InitManifestFolder(repoUrls string) bool {
	sshKeyPath := os.Getenv("SSH_KEY_PATH")
	g := git.NewGit(sshKeyPath, "user", "mail")

	urls := strings.Split(repoUrls, ",")
	for _, url := range urls {
		g.Clone(url, ReposBaseFolder, "", true)
	}

	return true
}

func GetImagesByRepo(folder string, repoUrl string, extensions []string) ([]string, []models.Image, map[string]models.ManifestBase, error) {
	pattern := ".*(" + repoUrl + ".*)"
	replacement := "${1}"
	values, images, projects, err := GetImages(folder, pattern, replacement, extensions, true)
	if err != nil {
		fmt.Println(err)
	}
	return values, images, projects, err
}

func GetImages(folder string, pattern string, replacement string, extensions []string, recursive bool) ([]string, []models.Image, map[string]models.ManifestBase, error) {
	values, err := filesystem.GetRegexInAllFiles(folder, true, pattern, replacement, extensions)
	if err != nil {
		fmt.Println(err)
	}

	images := []models.Image{}
	for _, value := range values {
		image := GetImage(value)
		// fmt.Println(image.Name)
		images = append(images, image)
		// break
	}

	projects := GetImageProjects(images)
	return values, images, projects, err
}

func GetImage(image string) models.Image {
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
	return models.Image{Tag: tag, Name: name, Project: project, BaseUrl: baseUrl}
}

func GetImageProjects(images []models.Image) map[string]models.ManifestBase {
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

	return bases
}
