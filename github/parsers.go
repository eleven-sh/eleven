package github

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/eleven-sh/eleven/entities"
	giturls "github.com/whilp/git-urls"
)

type ParsedRepositoryName struct {
	Owner         string
	ExplicitOwner bool
	Name          string
}

func ParseRepositoryName(
	repositoryName string,
	defaultRepositoryOwner string,
) (*ParsedRepositoryName, error) {

	errInvalidGitHubURL := errors.New("ErrInvalidGitHubURL")

	if len(repositoryName) == 0 {
		return nil, errInvalidGitHubURL
	}

	// Handle "git@github.com:eleven-sh/eleven.git"
	repositoryNameAsURL, err := giturls.Parse(repositoryName)

	if err != nil {
		// Handle "https://github.com/eleven-sh/eleven.git"
		repositoryNameAsURL, err = url.Parse(repositoryName)
	}

	// Not an URL (eg: "eleven") or only path (eg: "eleven-sh/eleven")
	if err != nil || len(repositoryNameAsURL.Hostname()) == 0 {
		repositoryNameParts := strings.Split(repositoryName, "/")

		if len(repositoryNameParts) > 2 {
			return nil, errInvalidGitHubURL
		}

		if len(repositoryNameParts) == 1 { // "eleven"
			return &ParsedRepositoryName{
				ExplicitOwner: false,
				Owner:         defaultRepositoryOwner,
				Name:          repositoryNameParts[0],
			}, nil
		}

		// Starts or ends with "/"
		if len(repositoryNameParts[0]) == 0 ||
			len(repositoryNameParts[1]) == 0 {

			return nil, errInvalidGitHubURL
		}

		return &ParsedRepositoryName{ // "eleven-sh/eleven"
			ExplicitOwner: true,
			Owner:         repositoryNameParts[0],
			Name:          repositoryNameParts[1],
		}, nil
	}

	host := repositoryNameAsURL.Hostname()
	if host != "github.com" {
		return nil, errInvalidGitHubURL
	}

	path := strings.TrimPrefix(repositoryNameAsURL.Path, "/")
	pathComponents := strings.Split(path, "/")

	if len(pathComponents) < 2 {
		return nil, errInvalidGitHubURL
	}

	githubRepositoryOwner := pathComponents[0]
	githubRepositoryName := strings.TrimSuffix(pathComponents[1], ".git")

	return &ParsedRepositoryName{
		ExplicitOwner: true,
		Owner:         githubRepositoryOwner,
		Name:          githubRepositoryName,
	}, nil
}

func BuildGitHTTPURL(repoName *ParsedRepositoryName) entities.EnvRepositoryGitURL {
	return entities.EnvRepositoryGitURL(fmt.Sprintf(
		"https://github.com/%s/%s.git",
		url.PathEscape(repoName.Owner),
		url.PathEscape(repoName.Name),
	))
}

func BuildGitURL(repoName *ParsedRepositoryName) entities.EnvRepositoryGitURL {
	return entities.EnvRepositoryGitURL(fmt.Sprintf(
		"git@github.com:%s/%s.git",
		url.PathEscape(repoName.Owner),
		url.PathEscape(repoName.Name),
	))
}
