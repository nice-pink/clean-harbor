package main

import (
	"os"

	"github.com/nice-pink/clean-harbor/pkg/cleaner"
	"github.com/nice-pink/clean-harbor/pkg/harbor"
	"github.com/nice-pink/clean-harbor/pkg/models"
	"github.com/nice-pink/clean-harbor/pkg/request"
)

func main() {
	config := harbor.HarborConfig{
		DryRun:    true,
		HarborUrl: os.Getenv("HARBOR_API"),
		BasicAuth: models.Auth{BasicUser: os.Getenv("HARBOR_USERNAME"), BasicPassword: os.Getenv("HARBOR_PASSWORD")},
	}

	requester := &request.Requester{}
	h := harbor.NewHarbor(requester, config)

	cleaner := cleaner.NewCleaner(h)
	extensions := []string{".yaml"}
	cleaner.FindUnused("../../pkg/test/repo", "quay.io", extensions)
}
