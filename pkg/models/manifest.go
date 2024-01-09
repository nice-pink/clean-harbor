package models

type ManifestBase struct {
	Name     string
	Projects map[string]ManifestProject
}

type ManifestProject struct {
	Name  string
	Repos map[string]ManifestRepo
}

type ManifestRepo struct {
	Name string
	Tags []string
}
