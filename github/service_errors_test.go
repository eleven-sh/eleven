package github

import (
	"errors"
	"net/http"
	"testing"

	"github.com/google/go-github/v43/github"
)

func TestServiceIsNotFoundError(t *testing.T) {
	testCases := []struct {
		test           string
		passedError    error
		expectedReturn bool
	}{
		{
			test:           "with base error",
			passedError:    errors.New(""),
			expectedReturn: false,
		},

		{
			test:           "with empty github error response",
			passedError:    &github.ErrorResponse{},
			expectedReturn: false,
		},

		{
			test: "with bad github error response status code",
			passedError: &github.ErrorResponse{
				Response: &http.Response{
					StatusCode: 403,
				},
			},
			expectedReturn: false,
		},

		{
			test: "with valid github error response",
			passedError: &github.ErrorResponse{
				Response: &http.Response{
					StatusCode: 404,
				},
			},
			expectedReturn: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			service := NewService()

			ret := service.IsNotFoundError(tc.passedError)

			if ret != tc.expectedReturn {
				t.Fatalf(
					"expected return to equal '%v', got '%v'",
					tc.expectedReturn,
					ret,
				)
			}
		})
	}
}

func TestServiceIsInvalidAccessTokenError(t *testing.T) {
	testCases := []struct {
		test           string
		passedError    error
		expectedReturn bool
	}{
		{
			test:           "with base error",
			passedError:    errors.New(""),
			expectedReturn: false,
		},

		{
			test:           "with empty github error response",
			passedError:    &github.ErrorResponse{},
			expectedReturn: false,
		},

		{
			test: "with bad github error response status code",
			passedError: &github.ErrorResponse{
				Response: &http.Response{
					StatusCode: 403,
				},
			},
			expectedReturn: false,
		},

		{
			test: "with valid github error response",
			passedError: &github.ErrorResponse{
				Response: &http.Response{
					StatusCode: 401,
				},
			},
			expectedReturn: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			service := NewService()

			ret := service.IsInvalidAccessTokenError(tc.passedError)

			if ret != tc.expectedReturn {
				t.Fatalf(
					"expected return to equal '%v', got '%v'",
					tc.expectedReturn,
					ret,
				)
			}
		})
	}
}
