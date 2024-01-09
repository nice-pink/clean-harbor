package models

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
