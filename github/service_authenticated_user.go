package github

import (
	"context"
	"errors"

	"github.com/google/go-github/v43/github"
)

type AuthenticatedUser struct {
	PrimaryEmail string
	Username     string
	FullName     string
}

func (s Service) GetAuthenticatedUser(
	accessToken string,
) (*AuthenticatedUser, error) {

	client := s.buildClient(accessToken)

	var primaryEmail string
	var getPrimaryEmailErr error
	var getPrimaryEmailChan = make(chan struct{})

	go func() {
		primaryEmail, getPrimaryEmailErr = s.getAuthenticatedUserPrimaryEmail(accessToken)

		close(getPrimaryEmailChan)
	}()

	// Passing the empty string will
	// fetch the authenticated user
	userName := ""
	user, _, getUserErr := client.Users.Get(context.TODO(), userName)

	<-getPrimaryEmailChan

	if getUserErr != nil {
		return nil, getUserErr
	}

	if getPrimaryEmailErr != nil {
		return nil, getPrimaryEmailErr
	}

	return &AuthenticatedUser{
		Username:     user.GetLogin(),
		FullName:     user.GetName(),
		PrimaryEmail: primaryEmail,
	}, nil
}

func (s Service) getAuthenticatedUserPrimaryEmail(
	accessToken string,
) (string, error) {

	client := s.buildClient(accessToken)

	emails, _, err := client.Users.ListEmails(context.TODO(), nil)

	if err != nil {
		return "", err
	}

	return getPrimaryEmail(emails)
}

func getPrimaryEmail(emails []*github.UserEmail) (string, error) {
	for _, email := range emails {
		if email.GetPrimary() {
			return email.GetEmail(), nil
		}
	}

	return "", errors.New("no primary email found")
}
