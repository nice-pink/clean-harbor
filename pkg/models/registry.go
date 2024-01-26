package models

type RegistryBase struct {
	Name     string
	Projects map[string]RegistryProject
}

type RegistryProject struct {
	Name      string
	Id        int
	RepoCount int
	Repos     map[string]RegistryRepo
}

type RegistryRepo struct {
	Name          string
	Id            int
	ArtifactCount int
	Artifacts     []RegistryArtifact
}

type RegistryTag struct {
	Name    string
	Created string
}

type RegistryArtifact struct {
	Tags    []RegistryTag
	ID      int
	Digest  string
	Created string
}
