package github

import (
	"errors"
	"testing"

	"github.com/google/go-github/v43/github"
)

func TestGetPrimaryEmail(t *testing.T) {
	stringP := func(s string) *string {
		return &s
	}
	boolP := func(b bool) *bool {
		return &b
	}

	testCases := []struct {
		test                 string
		emails               []*github.UserEmail
		expectedPrimaryEmail string
		expectedError        error
	}{
		{
			test: "with primary email",
			emails: []*github.UserEmail{
				{
					Email:   stringP("test@test.test"),
					Primary: boolP(true),
				},

				{
					Email:   stringP("test2@test.test"),
					Primary: boolP(false),
				},
			},
			expectedPrimaryEmail: "test@test.test",
			expectedError:        nil,
		},

		{
			test: "without primary email",
			emails: []*github.UserEmail{
				{
					Email:   stringP("test@test.test"),
					Primary: boolP(false),
				},

				{
					Email:   stringP("test2@test.test"),
					Primary: boolP(false),
				},
			},
			expectedError: errors.New(""),
		},

		{
			test:          "with empty emails",
			emails:        []*github.UserEmail{},
			expectedError: errors.New(""),
		},

		{
			test:          "with nil emails",
			emails:        nil,
			expectedError: errors.New(""),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			primaryEmail, err := getPrimaryEmail(
				tc.emails,
			)

			if tc.expectedError != nil && err == nil {
				t.Fatalf("expected error, got nothing")
			}

			if err != nil && tc.expectedError == nil {
				t.Fatalf("expected no error, got '%+v'", err)
			}

			if tc.expectedError != nil {
				return
			}

			if tc.expectedPrimaryEmail != primaryEmail {
				t.Fatalf(
					"expected primary email to equal '%s', got '%s'",
					tc.expectedPrimaryEmail,
					primaryEmail,
				)
			}
		})
	}
}
