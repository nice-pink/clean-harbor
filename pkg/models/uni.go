package models

import "fmt"

type UniBase struct {
	Name     string
	Projects []UniProject
}

type UniProject struct {
	Name  string
	Repos []UniRepo
}

type UniRepo struct {
	Name string
	Tags []string
}

func (u *UniBase) Print() {
	fmt.Println(u.Name)
	for _, project := range u.Projects {
		project.Print()
	}
}

func (u *UniProject) Print() {
	fmt.Println("  ", u.Name)
	for _, repo := range u.Repos {
		repo.Print()
	}
}

func (u *UniRepo) Print() {
	fmt.Println("    ", u.Name, "-", u.Tags)
}
