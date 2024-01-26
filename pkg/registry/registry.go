package registry

import (
	"github.com/nice-pink/clean-harbor/pkg/models"
)

func BuildUniModels(projects map[string]models.RegistryProject, baseUrl string) []models.UniBase {
	uBases := []models.UniBase{}
	uProjects := []models.UniProject{}
	// fmt.Println("Base:", base.Name)
	for _, project := range projects {
		uRepos := []models.UniRepo{}
		// fmt.Println("	Project:", project.Name)
		// fmt.Println("		", project.Name, "has repos", strconv.Itoa(len(project.Repos)))
		for _, repo := range project.Repos {
			// fmt.Println("		", repo.Name, repo.Tags)
			tags := []models.UniTag{}
			for _, artifact := range repo.Artifacts {
				if len(artifact.Tags) == 0 {
					tags = append(tags, models.UniTag{Name: "", Digest: artifact.Digest})
				}
				for _, tag := range artifact.Tags {
					tags = append(tags, models.UniTag{Name: tag.Name, Digest: artifact.Digest})
				}
			}
			uRepos = append(uRepos, models.UniRepo{Name: repo.Name, Tags: tags})
		}
		uProjects = append(uProjects, models.UniProject{Name: project.Name, Repos: uRepos})
	}
	uBases = append(uBases, models.UniBase{Name: baseUrl, Projects: uProjects})
	return uBases
}
