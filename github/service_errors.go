package github

import "github.com/google/go-github/v43/github"

func (s Service) IsNotFoundError(err error) bool {
	if githubErr, ok := err.(*github.ErrorResponse); ok &&
		githubErr.Response != nil &&
		githubErr.Response.StatusCode == 404 {

		return true
	}

	return false
}

func (s Service) IsInvalidAccessTokenError(err error) bool {
	if githubErr, ok := err.(*github.ErrorResponse); ok &&
		githubErr.Response != nil &&
		githubErr.Response.StatusCode == 401 {

		return true
	}

	return false
}
