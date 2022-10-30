package github

import (
	"context"

	gogithub "github.com/google/go-github/v43/github"
	"golang.org/x/oauth2"
)

type Service struct{}

func NewService() Service {
	return Service{}
}

func (s Service) buildClient(accessToken string) *gogithub.Client {
	oAuthTokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: accessToken,
		},
	)

	oAuthClient := oauth2.NewClient(
		context.TODO(),
		oAuthTokenSource,
	)

	return gogithub.NewClient(oAuthClient)
}
