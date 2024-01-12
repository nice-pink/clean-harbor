package cleaner

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/nice-pink/clean-harbor/pkg/harbor"
	"github.com/nice-pink/clean-harbor/pkg/manifestcrawler"
	"github.com/nice-pink/clean-harbor/pkg/models"
	"github.com/nice-pink/goutil/pkg/log"
)

// interfaces

type Harbor interface {
	GetAll() map[string]models.HarborProject
}

// cleaner

type Cleaner struct {
	h            Harbor
	dryRun       bool
	TAGS_HISTORY int
}

func NewCleaner(harbor Harbor, dryRun bool, tagsHistory int) *Cleaner {
	return &Cleaner{
		h:            harbor,
		dryRun:       dryRun,
		TAGS_HISTORY: tagsHistory,
	}
}

// clean unused

func (c *Cleaner) Remove(models []models.UniBase) (failed []string, succeed []string) {
	failed = append(failed, "one")

	if c.dryRun {
		fmt.Println("Dry run.")
	}

	return failed, nil
}

// find unused

func (c *Cleaner) FindUnused(repoFolder string, baseUrl string, extensions []string, ignoreUnsuedProjects bool, ignoreUnsuedRepos bool) []models.UniBase {
	// unused := []models.UniBase{}
	// unused = append(unused, models.UniBase{})

	// get harbor and manifest models
	harborModels, harborProjects, manifestModels := c.generateModels(repoFolder, baseUrl, extensions)
	unused := harborModels

	// get base project
	harborUniProjects := harborModels[0]

	fmt.Println("models")

	// find unused
	if _, ok := manifestModels[baseUrl]; ok {
		fmt.Println("has base", baseUrl)
		projects := []models.UniProject{}
		// get projects
		for _, hProject := range harborUniProjects.Projects {
			fmt.Print("project: '", hProject.Name, "'")
			if _, ok := manifestModels[baseUrl].Projects[hProject.Name]; ok {
				fmt.Println(" IS known! ✅")
				repos := []models.UniRepo{}
				// get repos
				for _, hRepo := range hProject.Repos {
					fmt.Print("  repo: '", hRepo.Name, "'")
					if mRepo, ok := manifestModels[baseUrl].Projects[hProject.Name].Repos[hRepo.Name]; ok {
						fmt.Println(" IS known! ✅")
						// get unused tags
						// unused[0].Projects[pIndex].Repos[rIndex].Tags = c.getUnusedTags(hRepo.Tags, mRepo.Tags)
						unusedTags := c.getUnusedTags(hRepo.Tags, mRepo.Tags)
						if !ignoreUnsuedRepos || len(unusedTags) > 0 {
							hRepo.Tags = unusedTags
							repos = append(repos, hRepo)
						} else {
							fmt.Println(" no tags to remove")
						}
					} else {
						fmt.Println(" UNUSED! 💥")
						if !ignoreUnsuedRepos {
							repos = append(repos, hRepo)
						}
					}
				}
				log.Info("append?", strconv.Itoa(len(repos)))
				hProject.Repos = repos

				if len(repos) > 0 || !ignoreUnsuedProjects {
					log.Info("append", strconv.Itoa(len(repos)))
					projects = append(projects, hProject)
				}
			} else {
				// unknown project
				fmt.Println(" UNUSED! 💥")
				if !ignoreUnsuedProjects {
					log.Info("append")
					projects = append(projects, hProject)
				}
			}
		}
		unused[0].Projects = projects
	}

	// log.Info("\n\nUnsued:")
	// for _, base := range unused {
	// 	base.Print()
	// }
	// log.Info()

	log.Info()
	log.Info("------------------------")
	log.Info("Unused artifacts:")
	unusedArtifacts := c.getUnusedArtifacts(unused, harborProjects, baseUrl)
	// for _, artifact := range unusedArtifacts {
	// 	artifact.Print()
	// }
	printImages("bin/unused_artifacts.txt", unusedArtifacts)

	return unused
}

func (c *Cleaner) getUnusedArtifacts(unused []models.UniBase, harborProjects map[string]models.HarborProject, baseUrl string) (unusedArtifacts []models.Image) {
	unusedArtifacts = []models.Image{}

	unusedProjects := unused[0].Projects
	for _, project := range unusedProjects {
		for _, repo := range project.Repos {
			if len(repo.Tags) == 0 {
				continue
			}

			// find tags and get list of digests, which are used to reference artifacts in harbor.
			tag := repo.Tags[0]

			repoKey := project.Name + "/" + repo.Name
			hArtifacts := harborProjects[project.Name].Repos[repoKey].Artifacts
			index := IndexOfTag(hArtifacts, tag)
			if index > 0 {
				// if index < len(hArtifacts) {
				for _, artifact := range hArtifacts[index:] {
					image := models.Image{Name: repo.Name, Project: project.Name, Tag: artifact.Digest, BaseUrl: baseUrl}
					unusedArtifacts = append(unusedArtifacts, image)
				}
				// }
			}
		}
	}

	return unusedArtifacts
}

func (c *Cleaner) getUnusedTags(harborTags []string, manifestTags []string) []string {
	tags := harborTags
	countTags := len(tags)
	// use artifact tags to compare with tags
	indeces := []int{}
	knownTags := []string{}
	for _, mTag := range manifestTags {
		index := IndexOf(tags, mTag)
		if index >= 0 {
			indeces = append(indeces, index)
			knownTags = append(knownTags, mTag)
		} else {
			fmt.Println("WARNING!!! Tag does not exist in registry: ", mTag)
		}
	}
	fmt.Println("    num existing tags:", strconv.Itoa(countTags))
	fmt.Println("    used tags:", knownTags)

	if len(indeces) > 0 {
		// get indeces and max index
		fmt.Print("    tag indeces:", indeces)
		sort.Ints(indeces)
		maximum := indeces[len(indeces)-1]
		fmt.Println(" ---> max:", strconv.Itoa(maximum))

		// are there unused tags?
		threshold := maximum + c.TAGS_HISTORY
		if countTags > threshold {
			fmt.Println("    unused tags:", tags[threshold:])
			return tags[maximum+1:]
		}
	}

	return nil
}

// get models

func (c *Cleaner) generateModels(repoFolder string, baseUrl string, extensions []string) (harborUniModels []models.UniBase, harborProjects map[string]models.HarborProject, manifestModels map[string]models.ManifestBase) {
	// generate harbor models
	harborProjects = c.h.GetAll()
	if len(harborProjects) == 0 {
		fmt.Println("No harbor projects.")
		return
	}
	harborUniModels = harbor.BuildUniModels(harborProjects, baseUrl)

	// generate manifest models
	_, _, manifestProjects, err := manifestcrawler.GetImagesByRepo(repoFolder, baseUrl, extensions)
	if err != nil || len(manifestProjects) == 0 {
		fmt.Println("No harbor projects.")
		return
	}

	// return
	return harborUniModels, harborProjects, manifestProjects
}

func (c *Cleaner) generateUniModels(repoFolder string, baseUrl string, extensions []string) (harborUniModels []models.UniBase, manifestUniModels []models.UniBase) {
	// generate harbor models
	harborProjects := c.h.GetAll()
	if len(harborProjects) == 0 {
		fmt.Println("No harbor projects.")
		return
	}
	harborModels := harbor.BuildUniModels(harborProjects, baseUrl)

	// generate manifest models
	_, _, manifestProjects, err := manifestcrawler.GetImagesByRepo(repoFolder, baseUrl, extensions)
	if err != nil || len(manifestProjects) == 0 {
		fmt.Println("No harbor projects.")
		return
	}
	manifestModels := manifestcrawler.BuildUniModels(manifestProjects)

	// return
	return harborModels, manifestModels
}

// helper

func IndexOf(slice []string, value string) int {
	for p, v := range slice {
		if v == value {
			return p
		}
	}
	return -1
}

func IndexOfTag(artifacts []models.HarborArtifact, tag string) int {
	index := 0
	for _, artifact := range artifacts {
		for _, aTag := range artifact.Tags {
			if aTag.Name == tag {
				return index
			}
		}
		index++
	}
	return -1
}

func printImages(filePath string, values []models.Image) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// print line by line and write to file
	for _, value := range values {
		value.Print()
		fmt.Fprintln(f, value.ToString())
	}
	return nil
}
