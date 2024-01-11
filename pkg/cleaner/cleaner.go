package cleaner

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/nice-pink/clean-harbor/pkg/harbor"
	"github.com/nice-pink/clean-harbor/pkg/manifestcrawler"
	"github.com/nice-pink/clean-harbor/pkg/models"
)

// interfaces

type Harbor interface {
	GetAll() []models.HarborProject
}

// cleaner

type Cleaner struct {
	h            Harbor
	TAGS_HISTORY int
}

func NewCleaner(harbor Harbor, tagsHistory int) *Cleaner {
	return &Cleaner{
		h:            harbor,
		TAGS_HISTORY: tagsHistory,
	}
}

// clean unused

func (c *Cleaner) Remove(models []models.UniBase, dryRun bool) (failed []string, succeed []string) {
	failed = append(failed, "one")

	return failed, nil
}

// find unused

func (c *Cleaner) FindUnused(repoFolder string, baseUrl string, extensions []string) []models.UniBase {
	// unused := []models.UniBase{}
	// unused = append(unused, models.UniBase{})

	// get harbor and manifest models
	harborModels, manifestModels := c.generateModels(repoFolder, baseUrl, extensions)
	unused := harborModels

	// get base project
	harborProjects := harborModels[0]

	// find unused
	if _, ok := manifestModels[baseUrl]; ok {
		fmt.Println("has base", baseUrl)
		for pIndex, project := range harborProjects.Projects {
			fmt.Print("project: '", project.Name, "'")
			if _, ok := manifestModels[baseUrl].Projects[project.Name]; ok {
				fmt.Println(" IS known! ✅")
				for rIndex, hRepo := range project.Repos {
					fmt.Print("  repo: '", hRepo.Name, "'")
					if mRepo, ok := manifestModels[baseUrl].Projects[project.Name].Repos[hRepo.Name]; ok {
						fmt.Println(" IS known! ✅")
						// get unused tags
						unused[0].Projects[pIndex].Repos[rIndex].Tags = c.getUnusedTags(hRepo.Tags, mRepo.Tags)
					} else {
						fmt.Println(" UNUSED! 💥")
					}
				}
			} else {
				// unknown project
				fmt.Println(" UNUSED! 💥")
				// unused[0].Projects = append(unused[0].Projects, project)
			}
		}
	}

	fmt.Println("\n\nUnsued:")
	for _, base := range unused {
		base.Print()
	}

	return unused
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

	return nil
}

// get models

func (c *Cleaner) generateModels(repoFolder string, baseUrl string, extensions []string) (harborUniModels []models.UniBase, manifestModels map[string]models.ManifestBase) {
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

	// return
	return harborModels, manifestProjects
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
