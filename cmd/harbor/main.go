package main

import (
	"os"

	"github.com/nice-pink/clean-harbor/pkg/harbor"
)

func main() {
	config := harbor.HarborConfig{
		DryRun:    true,
		HarborUrl: os.Getenv("HARBOR_API"),
	}

	r := harbor.NewHttpClient(os.Getenv("HARBOR_USERNAME"), os.Getenv("HARBOR_PASSWORD"), 30, 1)
	h := harbor.NewHarbor(r, config)

	h.GetAll("")
	// h.GetProjects(1, 2, true)
}
