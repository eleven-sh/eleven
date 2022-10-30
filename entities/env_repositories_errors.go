package entities

type ErrEnvRepositoryNotFound struct {
	RepoOwner string
	RepoName  string
}

func (ErrEnvRepositoryNotFound) Error() string {
	return "ErrEnvRepositoryNotFound"
}

type ErrEnvDuplicatedRepositories struct {
	RepoOwner string
	RepoName  string
}

func (ErrEnvDuplicatedRepositories) Error() string {
	return "ErrEnvDuplicatedRepositories"
}
