package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/nice-pink/clean-harbor/pkg/cleaner"
	"github.com/nice-pink/clean-harbor/pkg/harbor"
	"github.com/nice-pink/clean-harbor/pkg/manifestcrawler"
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
	reposDestFolder := flag.String("reposDestFolder", "", "Repo base folder.")
	repoUrls := flag.String("repoUrls", "", "Comma separated list of repoUrls. E.g. git@github.com:nice-pink/goutil.git,git@github.com:nice-pink/clean-harbor.git")
	registryBase := flag.String("registryBase", "", "Registry base which is used to identify images. E.g. 'quay.io'")
	filterProjects := flag.String("filterProjects", "", "Comma separated list of projects to search for. All others are ignored. E.g. websites,services")
	ignoreUnusedProjects := flag.Bool("ignoreUnusedProjects", false, "Unused projects are ignored. Could be, because they are handled differently e.g. pull through cache.")
	ignoreUnusedRepos := flag.Bool("ignoreUnusedRepos", false, "Unused repo are ignored. Could be, because they are currently unused.")
	flag.Parse()

	// if *repoUrls == "" {
	// 	log.Error("Please specify parameter: -repoUrls")
	// 	os.Exit(2)
	// }

	if *registryBase == "" {
		*registryBase = os.Getenv("REGISTRY_BASE")
	}
	if *registryBase == "" {
		*reposDestFolder = os.Getenv("REPO_FOLDER")
	}

	run(*reposDestFolder, *repoUrls, *registryBase, *ignoreUnusedProjects, *ignoreUnusedRepos, *filterProjects)
}

func run(reposDestFolder string, repoUrls string, registryBase string, ignoreUnusedProjects bool, ignoreUnusedRepos bool, filterProjectsString string) {
	start := time.Now()
	fmt.Println("Start:", start.Format(time.RFC3339))

	// checkout repo
	manifestcrawler.ReposBaseFolder = reposDestFolder
	if repoUrls != "" {
		log.Info("Checkout repos", repoUrls)
		manifestcrawler.InitManifestFolder(repoUrls)
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
		DryRun:    true,
		HarborUrl: os.Getenv("HARBOR_API"),
	}
	h := harbor.NewHarbor(r, config)

	// setup cleaner
	TAGS_HISTORY := 5
	dryRun := true
	cleaner := cleaner.NewCleaner(h, dryRun, TAGS_HISTORY)

	// get projects to filter by
	filterProjects := []string{}
	if filterProjectsString != "" {
		filterProjects = strings.Split(filterProjectsString, ",")
		log.Info("Only search for projects:", filterProjects)
	}

	// get unused projects
	extensions := []string{".yaml"}

	_, unused := cleaner.FindUnused(reposDestFolder, registryBase, extensions, filterProjects, ignoreUnusedProjects, ignoreUnusedRepos)
	npjson.DumpJson(unused, "bin/unused.json")

	// log duration
	end := time.Now()
	fmt.Println("End:", end.Format(time.RFC3339))
	fmt.Println("Duration:")
	duration := end.Sub(start)
	fmt.Println(duration)
}
