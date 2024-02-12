package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nice-pink/clean-harbor/pkg/cleaner"
	"github.com/nice-pink/clean-harbor/pkg/config"
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
// 	filterRepos := "feature-"
// 	cleaner.FindUnused("pkg/test/repo", "repo.url", extensions, filterProjects, filterRepos, false, false)
// }

func main() {
	baseFolder := flag.String("baseFolder", "bin", "Base folder.")
	reposDestFolder := flag.String("reposDestFolder", "repo", "Repo base folder.")
	repoUrls := flag.String("repoUrls", "", "Comma separated list of repoUrls. E.g. 'git@github.com:nice-pink/goutil.git,git@github.com:nice-pink/clean-harbor.git'")
	registryBase := flag.String("registryBase", "", "Registry base which is used to identify images. E.g. 'quay.io'")
	filterProjects := flag.String("filterProjects", "", "Comma separated list of projects to search for. All others are ignored. E.g. 'websites,services'")
	filterRepos := flag.String("filterRepos", "", "Search string contained in repo name. All others are ignored. E.g. '-feature-'")
	includeUnknownProjects := flag.Bool("includeUnknownProjects", false, "Unknown projects are included (and deleted). Be cautious: Could be unknown, because they are handled differently e.g. pull through cache.")
	includeUnknownRepos := flag.Bool("includeUnknownRepos", false, "Unknown repo are included (and deleted). Could be, because they are currently unused.")
	tagsHistory := flag.Int("tagsHistory", 5, "How many tags more than the oldest in use should be kept? Default=5")
	unknownHistory := flag.Int("unknownHistory", 0, "How many tags of unknown repos should be kept? Default=0")
	requestTimeout := flag.Int("requestTimeout", 30, "Timeout for any api request. Default=30")
	delete := flag.Bool("delete", false, "Should artifacts be deleted! This can't be undone!")
	dryRun := flag.Bool("dryRun", false, "Do dry run!")
	configPath := flag.String("configPath", "", "Config file incl. registry credentials.")
	// unusedArtifactsFilepath := flag.String("unusedArtifactsFilepath", "", "Set file path if only delete already found artifacts.")
	flag.Parse()

	// if *repoUrls == "" {
	// 	log.Error("Please specify parameter: -repoUrls")
	// 	os.Exit(2)
	// }

	// get config
	var cfg config.Config
	if *configPath != "" {
		cfg = config.GetConfig(*configPath)
		cfg.IsInitialised = true
		if cfg.IsInitialised {
			log.Info("Config file loaded.")
		}
	}

	if *dryRun {
		log.Info()
		log.Info("This is a Dry Run!")
		log.Info()
	}

	if *registryBase == "" {
		if cfg.IsInitialised && cfg.Registry.Base != "" {
			*registryBase = cfg.Registry.Base
		} else {
			*registryBase = os.Getenv("REGISTRY_BASE")
		}
	}
	if *registryBase == "" {
		if cfg.IsInitialised && cfg.Registry.Folder != "" {
			*reposDestFolder = cfg.Registry.Folder
		} else {
			*reposDestFolder = os.Getenv("REPO_FOLDER")
		}
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
		log.Info("YES")
	}

	// setup requester
	registryUser := ""
	if cfg.IsInitialised && cfg.Registry.User != "" {
		registryUser = cfg.Registry.User
	} else {
		registryUser = os.Getenv("REGISTRY_USERNAME")
	}
	registryPassword := ""
	if cfg.IsInitialised && cfg.Registry.Password != "" {
		registryPassword = cfg.Registry.Password
	} else {
		registryPassword = os.Getenv("REGISTRY_PASSWORD")
	}
	requestConfig := network.RequestConfig{
		Auth: network.Auth{
			BasicUser:     registryUser,
			BasicPassword: registryPassword,
		},
		Timeout: time.Duration(*requestTimeout),
	}
	r := network.NewRequester(requestConfig)

	// setup harbor
	registryApi := ""
	if cfg.IsInitialised && cfg.Registry.Api != "" {
		registryApi = cfg.Registry.Api
	} else {
		registryApi = os.Getenv("REGISTRY_API")
	}
	harborConfig := harbor.HarborConfig{
		DryRun:    *dryRun,
		HarborUrl: registryApi,
	}
	h := harbor.NewHarbor(r, harborConfig)

	// setup cleaner
	c := cleaner.NewCleaner(h, *dryRun, *tagsHistory, *unknownHistory)

	unusedArtifacts, unusedRepos := find(c, *baseFolder, repoDestFolder, *registryBase, !*includeUnknownProjects, !*includeUnknownRepos, *filterProjects, *filterRepos)

	if *delete {
		deleteUnused(c, unusedArtifacts, unusedRepos)
	}
}

func find(c *cleaner.Cleaner, baseFolder string, reposDestFolder string, registryBase string, ignoreUnusedProjects bool, ignoreUnusedRepos bool, filterProjectsString string, filterRepos string) (artifacts []models.Image, repos []models.Image) {
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

	artifacts, repos, unused := c.FindUnused(reposDestFolder, registryBase, extensions, filterProjects, filterRepos, ignoreUnusedProjects, ignoreUnusedRepos)

	// print unsued images
	unusedFilepath := filepath.Join(baseFolder, "unused.json")
	npjson.DumpJson(unused, unusedFilepath)

	// log.Info()
	log.Info("------------------------")
	log.Info("Unused artifacts:")
	log.Info(len(artifacts))
	unusedArtifactsFilepath := filepath.Join(baseFolder, "unused_artifacts.txt")
	cleaner.PrintImages(unusedArtifactsFilepath, artifacts, false)

	log.Info("------------------------")
	log.Info("Unused repos:")
	log.Info(len(repos))
	unusedReposFilepath := filepath.Join(baseFolder, "unused_repos.txt")
	cleaner.PrintImages(unusedReposFilepath, repos, false)

	// log duration
	end := time.Now()
	fmt.Println("End find:", end.Format(time.RFC3339))
	fmt.Println("Duration:")
	duration := end.Sub(start)
	fmt.Println(duration)

	return artifacts, repos
}

func deleteUnused(c *cleaner.Cleaner, artifacts []models.Image, repos []models.Image) map[string]error {
	// delete artifacts
	log.Info()
	log.Info("-------------------")
	log.Info("Delete artifacts:")
	start := time.Now()
	fmt.Println("Start delete:", start.Format(time.RFC3339))
	artifactErrors := c.Delete(artifacts)

	// delete repos
	log.Info()
	log.Info("-------------------")
	log.Info("Delete repos:")
	repoErrors := c.Delete(repos)
	// log duration
	end := time.Now()
	fmt.Println("End delete:", end.Format(time.RFC3339))
	fmt.Println("Duration:")
	duration := end.Sub(start)
	fmt.Println(duration)

	// merge and return
	MergeErrorMaps(artifactErrors, repoErrors)
	return artifactErrors
}

func MergeErrorMaps(e1 map[string]error, e2 map[string]error) {
	// no repo errors
	if e2 == nil {
		return
	}
	// no artifact errors
	if e1 == nil {
		e1 = e2
		return
	}
	// merge
	for k, v := range e2 {
		e1[k] = v
	}
}
