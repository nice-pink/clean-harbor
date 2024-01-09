package main

import (
	"os"

	"github.com/nice-pink/clean-harbor/pkg/harbor"
)

func main() {
	config := harbor.HarborConfig{
		DryRun:         true,
		HarborUrl:      os.Getenv("HARBOR_API"),
		HarborUser:     os.Getenv("HARBOR_USERNAME"),
		HarborPassword: os.Getenv("HARBOR_PASSWORD"),
	}
	harbor.Configure(config)
	// harbor.GetRepo("fluxmusic-builder", "websites")
	// harbor.GetRepos("websites", 1, 3)

	harbor.GetAll()
}
