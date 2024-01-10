package main

import (
	"os"

	"github.com/nice-pink/clean-harbor/pkg/harbor"
	"github.com/nice-pink/clean-harbor/pkg/network"
	"github.com/nice-pink/clean-harbor/pkg/request"
)

func main() {
	config := harbor.HarborConfig{
		DryRun:    true,
		HarborUrl: os.Getenv("HARBOR_API"),
		BasicAuth: network.Auth{BasicUser: os.Getenv("HARBOR_USERNAME"), BasicPassword: os.Getenv("HARBOR_PASSWORD")},
	}

	r := &request.Requester{}
	h := harbor.NewHarbor(r, config)

	// harbor.Configure(config)
	// // harbor.GetRepo("fluxmusic-builder", "websites")
	// // harbor.GetRepos("websites", 1, 3)

	h.GetAll()
}
