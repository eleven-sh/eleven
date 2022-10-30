package github

import (
	"context"
)

func (s Service) DoesRepositoryExist(
	accessToken string,
	repositoryOwner string,
	repositoryName string,
) (bool, error) {

	client := s.buildClient(accessToken)

	repository, _, err := client.Repositories.Get(
		context.TODO(),
		repositoryOwner,
		repositoryName,
	)

	if s.IsNotFoundError(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return repository != nil, nil
}
