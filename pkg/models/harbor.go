package models

type HarborBase struct {
	Name     string
	Projects map[string]HarborProject
}

type HarborProject struct {
	Name      string
	Id        int `json:"project_id"`
	RepoCount int `json:"repo_count"`
	Repos     map[string]HarborRepo
}

type HarborRepo struct {
	Name          string
	Id            int
	ArtifactCount int `json:"artifact_count"`
	Artifacts     []HarborArtifact
}

type HarborTag struct {
	Name    string
	Created string `json:"push_time"`
}

type HarborArtifact struct {
	Tags    []HarborTag
	ID      int
	Digest  string
	Created string `json:"push_time"`
}
