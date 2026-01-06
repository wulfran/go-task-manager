package requests

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCommonsSetFailed(t *testing.T) {
	var tests = []struct {
		name               string
		vr                 ValidationResult
		expectedValidation ValidationResult
		message            string
	}{
		{
			"first fail, clean message",
			ValidationResult{
				Validated: true,
			},
			ValidationResult{
				Validated: false,
				Message:   "lorem ipsum",
			},
			"lorem ipsum",
		},
		{
			"already failed, append message",
			ValidationResult{
				Validated: false,
				Message:   "lorem ipsum",
			},
			ValidationResult{
				Validated: false,
				Message:   "lorem ipsum, dolor et",
			},
			"dolor et",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.vr.SetFailed(tc.message)

			if diff := cmp.Diff(tc.vr, tc.expectedValidation); diff != "" {
				t.Errorf("invalid validation result state <-want, +got>\n%s", diff)
			}
		})
	}
}
