package main

import (
	"fmt"
	"os"
	"time"

	"github.com/nice-pink/clean-harbor/pkg/cleaner"
	"github.com/nice-pink/clean-harbor/pkg/harbor"
	npjson "github.com/nice-pink/goutil/pkg/json"
	"github.com/nice-pink/goutil/pkg/network"
)

// test
// func main() {
// 	h := &mock.MockHarbor{}
// 	TAGS_HISTORY := 5
//	dryRun := true

// 	cleaner := cleaner.NewCleaner(h, dryRun, TAGS_HISTORY)
// 	extensions := []string{".yaml"}
// 	cleaner.FindUnused("pkg/test/repo", "repo.url", extensions)
// }

func main() {
	start := time.Now()
	fmt.Println("Start:", start.Format(time.RFC3339))

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
	unused := cleaner.FindUnused(os.Getenv("REPO_FOLDER"), os.Getenv("REPO_BASE"), extensions)
	npjson.DumpJson(unused, "bin/unused.json")

	end := time.Now()
	fmt.Println("End:", end.Format(time.RFC3339))
	fmt.Println("Duration:")
	duration := end.Sub(start)
	fmt.Println(duration)
}
