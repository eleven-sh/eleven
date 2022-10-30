package entities

type EnvRepositoryGitURL string

type EnvRepository struct {
	Name          string              `json:"name"`
	Owner         string              `json:"owner"`
	ExplicitOwner bool                `json:"explicit_owner"`
	GitURL        EnvRepositoryGitURL `json:"git_url"`
	GitHTTPURL    EnvRepositoryGitURL `json:"git_http_url"`
}
