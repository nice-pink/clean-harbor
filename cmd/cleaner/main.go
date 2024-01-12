package main

import (
	"fmt"
	"os"
	"time"

	"github.com/nice-pink/clean-harbor/pkg/cleaner"
	"github.com/nice-pink/clean-harbor/pkg/harbor"
	"github.com/nice-pink/clean-harbor/pkg/manifestcrawler"
	"github.com/nice-pink/clean-harbor/pkg/test/mock"
	npjson "github.com/nice-pink/goutil/pkg/json"
	"github.com/nice-pink/goutil/pkg/log"
	"github.com/nice-pink/goutil/pkg/network"
)

// test
func main() {
	h := &mock.MockHarbor{}
	TAGS_HISTORY := 1
	dryRun := true

	cleaner := cleaner.NewCleaner(h, dryRun, TAGS_HISTORY)
	extensions := []string{}
	cleaner.FindUnused("pkg/test/repo", "repo.url", extensions, false, false)
}

// func main() {
// 	reposDestFolder := flag.String("reposDestFolder", "", "Repo base folder.")
// 	repoUrls := flag.String("repoUrls", "", "Commaseparated list of repoUrls. github.com/nice-pink/goutil,github.com/nice-pink/clean-harbor")
// 	registryBase := flag.String("registryBase", "", "Registry base which is used to identify images. E.g. 'quay.io'")
// 	ignoreUnusedProjects := flag.Bool("ignoreUnusedProjects", false, "Unused projects are ignored. Could be, because they are handled differently e.g. pull through cache.")
// 	ignoreUnusedRepos := flag.Bool("ignoreUnusedRepos", false, "Unused repo are ignored. Could be, because they are currently unused.")
// 	flag.Parse()

// 	// if *repoUrls == "" {
// 	// 	log.Error("Please specify parameter: -repoUrls")
// 	// 	os.Exit(2)
// 	// }

// 	if *registryBase == "" {
// 		*registryBase = os.Getenv("REGISTRY_BASE")
// 	}
// 	if *registryBase == "" {
// 		*reposDestFolder = os.Getenv("REPO_FOLDER")
// 	}

// 	run(*reposDestFolder, *repoUrls, *registryBase, *ignoreUnusedProjects, *ignoreUnusedRepos)
// }

func run(reposDestFolder string, repoUrls string, registryBase string, ignoreUnusedProjects bool, ignoreUnusedRepos bool) {
	start := time.Now()
	fmt.Println("Start:", start.Format(time.RFC3339))

	// checkout repo
	if repoUrls != "" {
		log.Info("Checkout repos", repoUrls)
		manifestcrawler.ReposBaseFolder = reposDestFolder
		manifestcrawler.InitManifestFolder(repoUrls)
	}

	requestConfig := network.RequestConfig{
		Auth: network.Auth{
			BasicUser:     os.Getenv("HARBOR_USERNAME"),
			BasicPassword: os.Getenv("HARBOR_PASSWORD"),
		},
	}

	config := harbor.HarborConfig{
		DryRun:    true,
		HarborUrl: os.Getenv("HARBOR_API"),
	}

	r := network.NewRequester(requestConfig)
	h := harbor.NewHarbor(r, config)
	TAGS_HISTORY := 5
	dryRun := true

	cleaner := cleaner.NewCleaner(h, dryRun, TAGS_HISTORY)
	extensions := []string{".yaml"}
	unused := cleaner.FindUnused(reposDestFolder, registryBase, extensions, ignoreUnusedProjects, ignoreUnusedRepos)
	npjson.DumpJson(unused, "bin/unused.json")

	end := time.Now()
	fmt.Println("End:", end.Format(time.RFC3339))
	fmt.Println("Duration:")
	duration := end.Sub(start)
	fmt.Println(duration)
}
