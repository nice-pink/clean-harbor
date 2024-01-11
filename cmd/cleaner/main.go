package main

import (
	"github.com/nice-pink/clean-harbor/pkg/cleaner"
	"github.com/nice-pink/clean-harbor/pkg/test/mock"
)

func main() {
	h := &mock.MockHarbor{}

	cleaner := cleaner.NewCleaner(h)
	extensions := []string{".yaml"}
	cleaner.FindUnused("pkg/test/repo", "repo.url", extensions)
}

// func main() {
// 	config := harbor.HarborConfig{
// 		DryRun:    true,
// 		HarborUrl: os.Getenv("HARBOR_API"),
// 		BasicAuth: models.Auth{BasicUser: os.Getenv("HARBOR_USERNAME"), BasicPassword: os.Getenv("HARBOR_PASSWORD")},
// 	}

// 	requester := &request.Requester{}
// 	h := harbor.NewHarbor(requester, config)

// 	cleaner := cleaner.NewCleaner(h)
// 	extensions := []string{".yaml"}
// 	cleaner.FindUnused("../../pkg/test/repo", "quay.io", extensions)
// }
