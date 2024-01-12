package main

import (
	"os"

	"github.com/nice-pink/clean-harbor/pkg/manifestcrawler"
	npjson "github.com/nice-pink/goutil/pkg/json"
)

func main() {
	base := os.Getenv("REPO_BASE")
	folder := os.Getenv("REPO_FOLDER")
	extensions := []string{".yaml", ".yml", ".kustomization"}
	_, images, projects, _ := manifestcrawler.GetImagesByRepo(folder, base, extensions)
	npjson.DumpJson(projects, "bin/repo.json")
	npjson.DumpJson(images, "bin/easy.json")

	// checkout repo
	// manifestcrawler.InitManifestFolder("git@github.com:nice-pink/random.git,git@github.com:nice-pink/goutil.git")
}
