package main

import (
	"os"

	"github.com/nice-pink/clean-harbor/pkg/manifestcrawler"
)

func main() {
	base := os.Getenv("REPO_BASE")
	folder := os.Getenv("REPO_FOLDER")
	extensions := []string{".yaml", ".yml", ".kustomization"}
	manifestcrawler.GetImagesByRepo(folder, base, extensions)
}
