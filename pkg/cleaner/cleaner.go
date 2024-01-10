package cleaner

import (
	"fmt"

	"github.com/nice-pink/clean-harbor/pkg/manifestcrawler"
	"github.com/nice-pink/clean-harbor/pkg/models"
)

// interfaces

type Harbor interface {
	GetAll() []models.HarborProject
}

// cleaner

type Cleaner struct {
	harbor Harbor
}

func NewCleaner(harbor Harbor) *Cleaner {
	return &Cleaner{
		harbor: harbor,
	}
}

// functions

func (c *Cleaner) FindUnused(repoFolder string, baseUrl string, extensions []string) {

	harborProjects := c.harbor.GetAll()
	if len(harborProjects) == 0 {
		fmt.Println("No harbor projects.")
		return
	}

	manifestProjects, err := manifestcrawler.GetImagesByRepo(repoFolder, baseUrl, extensions)
	if err != nil || len(manifestProjects) == 0 {
		fmt.Println("No harbor projects.")
		return
	}
}
