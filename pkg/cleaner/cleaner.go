package cleaner

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/nice-pink/clean-harbor/pkg/manifestcrawler"
	"github.com/nice-pink/clean-harbor/pkg/models"
	"github.com/nice-pink/clean-harbor/pkg/registry"
	"github.com/nice-pink/goutil/pkg/log"
)

// cleaner

type Cleaner struct {
	r              registry.Registry
	dryRun         bool // do not delete if dry run
	tagsHistory    int  // amount of artifacts kept for known repos additionally to the oldest known
	unknownHistory int  // amount of artifacts kept for unknown repos
}

func NewCleaner(registry registry.Registry, dryRun bool, tagsHistory int, unknownHistory int) *Cleaner {
	return &Cleaner{
		r:              registry,
		dryRun:         dryRun,
		tagsHistory:    tagsHistory,
		unknownHistory: unknownHistory,
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

func (c *Cleaner) FindUnused(repoFolder string, baseUrl string, extensions []string, filterProjects []string, filterRepos string, ignoreUnsuedProjects bool, ignoreUnsuedRepos bool) ([]models.Image, []models.Image, []models.UniBase) {
	// get harbor and manifest models
	regModels, regProjects, manifestModels := c.generateModels(repoFolder, baseUrl, extensions, filterProjects, filterRepos)
	unused := regModels

	// get base project
	regUniProjects := regModels[0]

	// find unused
	if _, ok := manifestModels[baseUrl]; ok {
		fmt.Println("has base", baseUrl)
		projects := []models.UniProject{}
		// get projects
		for _, hProject := range regUniProjects.Projects {
			fmt.Print("project: '", hProject.Name, "'")
			if _, ok := manifestModels[baseUrl].Projects[hProject.Name]; ok {
				fmt.Println(" IS known! ✅")
				repos := []models.UniRepo{}
				// get repos
				for _, hRepo := range hProject.Repos {
					if filterRepos != "" && !strings.Contains(hRepo.Name, filterRepos) {
						continue
					}
					fmt.Print("  repo: '", hRepo.Name, "'")

					mRepo, repoIsKnown := manifestModels[baseUrl].Projects[hProject.Name].Repos[hRepo.Name]
					if repoIsKnown {
						fmt.Println(" IS known! ✅")

						// get unused tags
						// unused[0].Projects[pIndex].Repos[rIndex].Tags = c.getUnusedTags(hRepo.Tags, mRepo.Tags)
						unusedTags := c.getUnusedTags(hRepo.Tags, mRepo.Tags, 0)
						if len(unusedTags) > 0 {
							hRepo.Tags = unusedTags
							repos = append(repos, hRepo)
						} else {
							fmt.Println(" no tags to remove")
						}
					} else {

						// delete repo (after artifacts were removed)
						fmt.Println(" UNUSED! 💥")
						if !ignoreUnsuedRepos {
							// log.Info("Include unused repo.")
							// keep
							if c.unknownHistory > 0 {
								unusedTags := c.getUnusedTags(hRepo.Tags, mRepo.Tags, c.unknownHistory)
								if len(unusedTags) > 0 {
									hRepo.Tags = unusedTags
								}
							}
							hRepo.Unused = true
							repos = append(repos, hRepo)
						}
					}
				}
				// log.Info("append?", strconv.Itoa(len(repos)))
				hProject.Repos = repos

				if len(repos) > 0 {
					// log.Info("append", strconv.Itoa(len(repos)))
					projects = append(projects, hProject)
				}
			} else {
				// unknown project
				fmt.Println(" UNUSED! 💥")
				if !ignoreUnsuedProjects {
					// log.Info("Include unused project.")
					hProject.Unused = true
					projects = append(projects, hProject)
				}
			}
		}
		unused[0].Projects = projects
	}

	// get unsed artifacts and repos
	unusedArtifacts, unsuedRepos := c.getUnusedItems(unused, regProjects, baseUrl)

	return unusedArtifacts, unsuedRepos, unused
}

// delete

func (c *Cleaner) Delete(images []models.Image) map[string]error {
	errors := map[string]error{}

	if !c.dryRun {
		log.Info("This is not a dry run!!!!")
	} else {
		log.Info("This is a DRY RUN!")
	}

	for _, image := range images {
		var err error
		if image.Tag == "" {
			log.Info("Delete Repo:", image.Name, image.Tag)
			if !c.dryRun {
				_, err = c.r.DeleteRepo(image.Project, image.Name)
			}
		} else {
			log.Info("Delete Artifact:", image.Name, image.Tag)
			if !c.dryRun {
				_, err = c.r.DeleteArtifact(image.Tag, image.Project, image.Name)
			}
		}

		if err != nil {

			key := image.Project + "/" + image.Name
			if image.Tag != "" {
				key += "/" + image.Tag
			}
			errors[key] = err
		}
	}

	return errors
}

//

func (c *Cleaner) getUnusedItems(unused []models.UniBase, regProjects map[string]models.RegistryProject, baseUrl string) (unusedArtifacts []models.Image, unusedRepos []models.Image) {
	unusedArtifacts = []models.Image{}
	unusedRepos = []models.Image{}

	unusedProjects := unused[0].Projects
	for _, project := range unusedProjects {
		for _, repo := range project.Repos {
			if repo.Unused || project.Unused {
				unusedRepos = append(unusedRepos, models.Image{Name: repo.Name, Project: project.Name, BaseUrl: baseUrl})
				// continue
			}

			if len(repo.Tags) == 0 {
				log.Warn("No tags for repo:", repo.Name)
				continue
			}

			// find tags and get list of digests, which are used to reference artifacts in harbor.
			tag := repo.Tags[0]

			hArtifacts := regProjects[project.Name].Repos[repo.Name].Artifacts

			index := IndexOfDigest(hArtifacts, tag.Digest)
			if index < 0 {
				continue
			}

			for _, artifact := range hArtifacts[index:] {
				image := models.Image{Name: repo.Name, Project: project.Name, Tag: artifact.Digest, BaseUrl: baseUrl}
				unusedArtifacts = append(unusedArtifacts, image)
			}
		}
	}

	return unusedArtifacts, unusedRepos
}

func (c *Cleaner) getUnusedTags(regTags []models.UniTag, manifestTags []string, minUsed int) []models.UniTag {
	// no tags in registry
	if len(regTags) == 0 {
		return nil
	}

	// get minimum
	if len(manifestTags) == 0 && minUsed > 0 {
		firstIndex := 0
		itemCounter := 0
		priorDigest := regTags[0].Digest
		for index, tag := range regTags {
			if tag.Digest != priorDigest {
				firstIndex = index
				priorDigest = tag.Digest
				itemCounter++
			}
			if itemCounter >= minUsed {
				break
			}
		}
		return regTags[firstIndex:]
	}

	tags := regTags
	countTags := len(tags)
	// use artifact tags to compare with tags
	indeces := []int{}
	knownTags := []string{}
	for _, mTag := range manifestTags {
		index := IndexInUniTags(tags, mTag)
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
		// fmt.Print("    tag indeces:", indeces) // value might be higher than expected, because one artifact can have multiple tags!
		sort.Ints(indeces)
		maximum := indeces[len(indeces)-1]
		// fmt.Println(" ---> max:", strconv.Itoa(maximum))

		// are there unused tags?
		threshold := maximum + c.tagsHistory + 1
		if countTags > threshold {
			fmt.Println("    UNUSED TAGS:", strconv.Itoa(len(tags[threshold:])), "starting from:", tags[threshold]) // tags[threshold:]
			return tags[threshold:]
		}
	}

	return nil
}

// get models

func (c *Cleaner) generateModels(repoFolder string, baseUrl string, extensions []string, filterProjects []string, filterRepos string) (regUniModels []models.UniBase, regProjects map[string]models.RegistryProject, manifestModels map[string]models.ManifestBase) {
	// generate harbor models
	if len(filterProjects) == 0 {
		regProjects = c.r.GetAll(filterRepos)
	} else {
		regProjects = map[string]models.RegistryProject{}
		for _, name := range filterProjects {
			log.Info("--- Get repos for project:")
			repos, _ := c.r.GetAllRepos(name, filterRepos, false)
			regProjects[name] = models.RegistryProject{Name: name, Repos: repos}
			log.Info("Got", strconv.Itoa(len(repos)), "repos.")
		}
		regProjects = c.r.EnrichReposWithArtificats(regProjects)
	}

	if len(regProjects) == 0 {
		fmt.Println("No harbor projects.")
		return
	}
	regUniModels = registry.BuildUniModels(regProjects, baseUrl)

	// generate manifest models
	_, _, manifestProjects, err := manifestcrawler.GetImagesByRepo(repoFolder, baseUrl, extensions)
	if err != nil || len(manifestProjects) == 0 {
		fmt.Println("No harbor projects.")
		return
	}

	// return
	return regUniModels, regProjects, manifestProjects
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

func IndexInUniTags(slice []models.UniTag, value string) int {
	for p, v := range slice {
		if v.Name == value {
			return p
		}
	}
	return -1
}

func IndexOfTag(artifacts []models.RegistryArtifact, tag string) int {
	for index, artifact := range artifacts {
		for _, aTag := range artifact.Tags {
			if aTag.Name == tag {
				return index
			}
		}
	}
	return -1
}

func IndexOfDigest(artifacts []models.RegistryArtifact, digest string) int {
	for index, artifact := range artifacts {
		if artifact.Digest == digest {
			return index
		}
	}
	return -1
}

func PrintImages(filePath string, values []models.Image, toStdout bool) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// print line by line and write to file
	for _, value := range values {
		if toStdout {
			value.Print()
		}

		fmt.Fprintln(f, value.ToString())
	}
	return nil
}
