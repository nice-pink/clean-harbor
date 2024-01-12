package main

import (
	"os"

	"github.com/nice-pink/clean-harbor/pkg/harbor"
	"github.com/nice-pink/goutil/pkg/network"
)

func main() {
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

	h.GetAll()
	// h.GetProjects(1, 2, true)
}
