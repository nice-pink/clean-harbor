package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nice-pink/clean-harbor/pkg/cleaner"
	"github.com/nice-pink/clean-harbor/pkg/harbor"
	"github.com/nice-pink/clean-harbor/pkg/manifestcrawler"
	"github.com/nice-pink/clean-harbor/pkg/models"
	"github.com/nice-pink/goutil/pkg/filesystem"
	npjson "github.com/nice-pink/goutil/pkg/json"
	"github.com/nice-pink/goutil/pkg/log"
	"github.com/nice-pink/goutil/pkg/network"
)

// test
// func main() {
// 	h := &mock.MockHarbor{}
// 	TAGS_HISTORY := 1
// 	dryRun := true

// 	cleaner := cleaner.NewCleaner(h, dryRun, TAGS_HISTORY)
// 	extensions := []string{}
// 	filterProjects := []string{"web"}
// 	cleaner.FindUnused("pkg/test/repo", "repo.url", extensions, filterProjects, false, false)
// }

func main() {
	action := flag.String("action", "find", "find, delete, findAndDelete")
	baseFolder := flag.String("baseFolder", "bin", "Base folder.")
	reposDestFolder := flag.String("reposDestFolder", "repo", "Repo base folder.")
	repoUrls := flag.String("repoUrls", "", "Comma separated list of repoUrls. E.g. git@github.com:nice-pink/goutil.git,git@github.com:nice-pink/clean-harbor.git")
	registryBase := flag.String("registryBase", "", "Registry base which is used to identify images. E.g. 'quay.io'")
	filterProjects := flag.String("filterProjects", "", "Comma separated list of projects to search for. All others are ignored. E.g. websites,services")
	ignoreUnusedProjects := flag.Bool("ignoreUnusedProjects", false, "Unused projects are ignored. Could be, because they are handled differently e.g. pull through cache.")
	ignoreUnusedRepos := flag.Bool("ignoreUnusedRepos", false, "Unused repo are ignored. Could be, because they are currently unused.")
	tagsHistory := flag.Int("tagsHistory", 5, "How many tags more than the oldest in use should be kept? Default=5")
	delete := flag.Bool("delete", false, "Should artifacts be deleted! This can't be undone!")
	// unusedArtifactsFilepath := flag.String("unusedArtifactsFilepath", "", "Set file path if only delete already found artifacts.")
	flag.Parse()

	// if *repoUrls == "" {
	// 	log.Error("Please specify parameter: -repoUrls")
	// 	os.Exit(2)
	// }

	DRY_RUN := !*delete

	if *registryBase == "" {
		*registryBase = os.Getenv("REGISTRY_BASE")
	}
	if *registryBase == "" {
		*reposDestFolder = os.Getenv("REPO_FOLDER")
	}

	// checkout repo
	repoDestFolder := filepath.Join(*baseFolder, *reposDestFolder)
	manifestcrawler.ReposBaseFolder = repoDestFolder
	if *repoUrls != "" {
		log.Info("Checkout repos", *repoUrls)
		manifestcrawler.InitManifestFolder(*repoUrls)
	} else {
		log.Info("Does repo folder exist?")
		if !filesystem.DirExists(repoDestFolder) {
			log.Error("Repo folder does not exist.")
			os.Exit(2)
		}
	}

	// setup requester
	requestConfig := network.RequestConfig{
		Auth: network.Auth{
			BasicUser:     os.Getenv("HARBOR_USERNAME"),
			BasicPassword: os.Getenv("HARBOR_PASSWORD"),
		},
	}
	r := network.NewRequester(requestConfig)

	// setup harbor
	config := harbor.HarborConfig{
		DryRun:    DRY_RUN,
		HarborUrl: os.Getenv("HARBOR_API"),
	}
	h := harbor.NewHarbor(r, config)

	// setup cleaner
	dryRun := DRY_RUN
	c := cleaner.NewCleaner(h, dryRun, *tagsHistory)

	unusedArtifacts := []models.Image{}

	if strings.ToUpper(*action) == "FIND" || strings.ToUpper(*action) == strings.ToUpper("findAndDelete") {
		_, unusedArtifacts = find(c, *baseFolder, repoDestFolder, *registryBase, *ignoreUnusedProjects, *ignoreUnusedRepos, *filterProjects, *tagsHistory)
	}
	//  else {
	// 	unusedFilepath = *unusedArtifactsFilepath
	// }

	if strings.ToUpper(*action) == "DELETE" || strings.ToUpper(*action) == strings.ToUpper("findAndDelete") {
		// if unusedFilepath == "" {
		// 	log.Error("No file path specified containing the artifacts to be deleted!")
		// 	log.Info("Either choose action 'find' or 'findAndDelete' OR specify -unusedArtifactsFilepath to an already existing file.")
		// 	os.Exit(2)
		// }
		log.Info()
		log.Info("-------------------")
		start := time.Now()
		fmt.Println("Start delete:", start.Format(time.RFC3339))
		c.Delete(unusedArtifacts)
		// log duration
		end := time.Now()
		fmt.Println("End delete:", end.Format(time.RFC3339))
		fmt.Println("Duration:")
		duration := end.Sub(start)
		fmt.Println(duration)
	}
}

func find(c *cleaner.Cleaner, baseFolder string, reposDestFolder string, registryBase string, ignoreUnusedProjects bool, ignoreUnusedRepos bool, filterProjectsString string, tagsHistory int) (string, []models.Image) {
	start := time.Now()
	fmt.Println("Start find:", start.Format(time.RFC3339))

	// get projects to filter by
	filterProjects := []string{}
	if filterProjectsString != "" {
		filterProjects = strings.Split(filterProjectsString, ",")
		log.Info("Only search for projects:", filterProjects)
	}

	// get unused projects
	extensions := []string{".yaml"}

	artifacts, unused := c.FindUnused(reposDestFolder, registryBase, extensions, filterProjects, ignoreUnusedProjects, ignoreUnusedRepos)

	// print unsued images
	unusedFilepath := filepath.Join(baseFolder, "unused.json")
	npjson.DumpJson(unused, unusedFilepath)

	// log.Info()
	log.Info("------------------------")
	log.Info("Unused artifacts:")
	log.Info(len(artifacts))
	unusedArtifactsFilepath := filepath.Join(baseFolder, "unused_artifacts.txt")
	cleaner.PrintImages(unusedArtifactsFilepath, artifacts, false)

	// log duration
	end := time.Now()
	fmt.Println("End find:", end.Format(time.RFC3339))
	fmt.Println("Duration:")
	duration := end.Sub(start)
	fmt.Println(duration)

	return unusedArtifactsFilepath, artifacts
}
