package github

import (
	"reflect"
	"testing"

	"github.com/eleven-sh/eleven/entities"
)

func TestParseRepositoryNameWithValidNames(t *testing.T) {
	testCases := []struct {
		test                   string
		repositoryName         string
		defaultRepositoryOwner string
		expectedParsedRepoName *ParsedRepositoryName
	}{
		{
			test:                   "with repository owner",
			repositoryName:         "eleven-sh/api",
			defaultRepositoryOwner: "jeremylevy",
			expectedParsedRepoName: &ParsedRepositoryName{
				Owner:         "eleven-sh",
				ExplicitOwner: true,
				Name:          "api",
			},
		},

		{
			test:                   "without repository owner",
			repositoryName:         "api",
			defaultRepositoryOwner: "jeremylevy",
			expectedParsedRepoName: &ParsedRepositoryName{
				Owner:         "jeremylevy",
				ExplicitOwner: false,
				Name:          "api",
			},
		},

		{
			test:                   "with Git URL",
			repositoryName:         "git@github.com:jeremylevy/fullstack-open.git",
			defaultRepositoryOwner: "jeremylevy",
			expectedParsedRepoName: &ParsedRepositoryName{
				Owner:         "jeremylevy",
				ExplicitOwner: true,
				Name:          "fullstack-open",
			},
		},

		{
			test:                   "with HTTP Git URL",
			repositoryName:         "https://github.com/jeremylevy/fullstack-open.git",
			defaultRepositoryOwner: "jeremylevy",
			expectedParsedRepoName: &ParsedRepositoryName{
				Owner:         "jeremylevy",
				ExplicitOwner: true,
				Name:          "fullstack-open",
			},
		},

		{
			test:                   "with HTTP root URL",
			repositoryName:         "https://github.com/foo-/UTwente-Usability",
			defaultRepositoryOwner: "jeremylevy",
			expectedParsedRepoName: &ParsedRepositoryName{
				Owner:         "foo-",
				ExplicitOwner: true,
				Name:          "UTwente-Usability",
			},
		},

		{
			test:                   "with HTTP non-root URL",
			repositoryName:         "https://github.com/foo-/UTwente-Usability/blob/master/modSDR.js",
			defaultRepositoryOwner: "jeremylevy",
			expectedParsedRepoName: &ParsedRepositoryName{
				Owner:         "foo-",
				ExplicitOwner: true,
				Name:          "UTwente-Usability",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			parsedRepoName, err := ParseRepositoryName(
				tc.repositoryName,
				tc.defaultRepositoryOwner,
			)

			if err != nil {
				t.Fatalf("expected no error, got '%+v'", err)
			}

			if !reflect.DeepEqual(tc.expectedParsedRepoName, parsedRepoName) {
				t.Fatalf(
					"expected parsed repository name to equal '%+v', got '%+v'",
					tc.expectedParsedRepoName,
					parsedRepoName,
				)
			}
		})
	}
}

func TestParseRepositoryNameWithInvalidNames(t *testing.T) {
	testCases := []struct {
		test                   string
		repositoryName         string
		defaultRepositoryOwner string
	}{
		{
			test:                   "with empty repository name",
			repositoryName:         "",
			defaultRepositoryOwner: "jeremylevy",
		},

		{
			test:                   "with invalid Git URL host",
			repositoryName:         "git@invalid.com:jeremylevy/fullstack-open.git",
			defaultRepositoryOwner: "jeremylevy",
		},

		{
			test:                   "with invalid HTTP Git URL host",
			repositoryName:         "https://invalid.com/jeremylevy/fullstack-open.git",
			defaultRepositoryOwner: "jeremylevy",
		},

		{
			test:                   "with invalid HTTP URL host",
			repositoryName:         "https://invalid.com/jeremylevy/fullstack-open/blob/master/index.js",
			defaultRepositoryOwner: "jeremylevy",
		},

		{
			test:                   "with starting slash",
			repositoryName:         "/index.php",
			defaultRepositoryOwner: "jeremylevy",
		},

		{
			test:                   "with ending slash",
			repositoryName:         "index.php/",
			defaultRepositoryOwner: "jeremylevy",
		},

		{
			test:                   "with starting and ending slash",
			repositoryName:         "/index.php/",
			defaultRepositoryOwner: "jeremylevy",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			_, err := ParseRepositoryName(
				tc.repositoryName,
				tc.defaultRepositoryOwner,
			)

			if err == nil {
				t.Fatalf("expected error got nothing")
			}
		})
	}
}

func TestBuildGitHTTPURL(t *testing.T) {
	testCases := []struct {
		test               string
		parsedRepoName     *ParsedRepositoryName
		expectedGitHTTPURL entities.EnvRepositoryGitURL
	}{
		{
			test: "with valid URL characters",
			parsedRepoName: &ParsedRepositoryName{
				Owner:         "foo-",
				ExplicitOwner: true,
				Name:          "api",
			},
			expectedGitHTTPURL: "https://github.com/foo-/api.git",
		},

		{
			test: "with invalid URL characters",
			parsedRepoName: &ParsedRepositoryName{
				Owner:         "eleven sh",
				ExplicitOwner: true,
				Name:          "repo/name",
			},
			expectedGitHTTPURL: "https://github.com/eleven%20sh/repo%2Fname.git",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			gitHTTPURL := BuildGitHTTPURL(
				tc.parsedRepoName,
			)

			if !reflect.DeepEqual(tc.expectedGitHTTPURL, gitHTTPURL) {
				t.Fatalf(
					"expected Git HTTP URL to equal '%s', got '%s'",
					tc.expectedGitHTTPURL,
					gitHTTPURL,
				)
			}
		})
	}
}

func TestBuildGitURL(t *testing.T) {
	testCases := []struct {
		test           string
		parsedRepoName *ParsedRepositoryName
		expectedGitURL entities.EnvRepositoryGitURL
	}{
		{
			test: "with valid URL characters",
			parsedRepoName: &ParsedRepositoryName{
				Owner:         "foo-",
				ExplicitOwner: true,
				Name:          "api",
			},
			expectedGitURL: "git@github.com:foo-/api.git",
		},

		{
			test: "with invalid URL characters",
			parsedRepoName: &ParsedRepositoryName{
				Owner:         "eleven sh",
				ExplicitOwner: true,
				Name:          "repo/name",
			},
			expectedGitURL: "git@github.com:eleven%20sh/repo%2Fname.git",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			gitURL := BuildGitURL(
				tc.parsedRepoName,
			)

			if !reflect.DeepEqual(tc.expectedGitURL, gitURL) {
				t.Fatalf(
					"expected Git URL to equal '%s', got '%s'",
					tc.expectedGitURL,
					gitURL,
				)
			}
		})
	}
}
