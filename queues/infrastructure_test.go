package queues

import (
	"errors"
	"testing"
)

func TestInfrastructureQueue(t *testing.T) {
	stringP := func(s string) *string {
		return &s
	}
	testCases := []struct {
		test                 string
		queue                InfrastructureQueue[*string]
		initialInfra         *string
		expectedInfra        *string
		expectedErrorMessage *string
	}{
		{
			test: "with no errors returned",
			queue: InfrastructureQueue[*string]{
				InfrastructureQueueSteps[*string]{
					func(infra *string) error {
						*infra += "a"
						return nil
					},

					func(infra *string) error {
						*infra += "a"
						return nil
					},
				},

				InfrastructureQueueSteps[*string]{
					func(infra *string) error {
						*infra += "b"
						return nil
					},

					func(infra *string) error {
						*infra += "b"
						return nil
					},

					func(infra *string) error {
						*infra += "b"
						return nil
					},
				},

				InfrastructureQueueSteps[*string]{
					func(infra *string) error {
						*infra += "c"
						return nil
					},
				},
			},
			initialInfra:         stringP(""),
			expectedInfra:        stringP("aabbbc"),
			expectedErrorMessage: nil,
		},

		{
			test: "with multiple errors returned",
			queue: InfrastructureQueue[*string]{
				InfrastructureQueueSteps[*string]{
					func(infra *string) error {
						*infra += "a"
						return nil
					},

					func(infra *string) error {
						return errors.New("my-error-1")
					},
				},

				InfrastructureQueueSteps[*string]{
					func(infra *string) error {
						*infra += "b"
						return nil
					},

					func(infra *string) error {
						*infra += "b"
						return errors.New("my-error-2")
					},
				},
			},
			initialInfra:         stringP(""),
			expectedInfra:        stringP("a"),
			expectedErrorMessage: stringP("my-error-1"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			err := tc.queue.Run(tc.initialInfra)

			if tc.expectedErrorMessage != nil && err == nil {
				t.Fatalf("expected error, got nothing")
			}

			if tc.expectedErrorMessage != nil &&
				err.Error() != *tc.expectedErrorMessage {

				t.Fatalf(
					"expected error message to equal '%s', got '%s'",
					*tc.expectedErrorMessage,
					err.Error(),
				)
			}

			if tc.expectedInfra != nil &&
				*tc.initialInfra != *tc.expectedInfra {

				t.Fatalf(
					"expected infrastructure to equal '%s', got '%s'",
					*tc.expectedInfra,
					*tc.initialInfra,
				)
			}
		})
	}
}
