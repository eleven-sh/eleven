package github

import (
	"context"

	"github.com/google/go-github/v43/github"
)

func (s Service) CreateSSHKey(
	accessToken string,
	keyPairName string,
	publicKeyContent string,
) (*github.Key, error) {

	client := s.buildClient(accessToken)

	key, _, err := client.Users.CreateKey(
		context.TODO(),
		&github.Key{
			Title: &keyPairName,
			Key:   &publicKeyContent,
		},
	)

	return key, err
}

func (s Service) RemoveSSHKey(
	accessToken string,
	sshKeyID int64,
) error {

	client := s.buildClient(accessToken)

	_, err := client.Users.DeleteKey(
		context.TODO(),
		sshKeyID,
	)

	return err
}
