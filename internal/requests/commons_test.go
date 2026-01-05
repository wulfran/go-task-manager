package requests

import (
	"testing"
)

func TestCommonsSetFailed(t *testing.T) {
	vr := ValidationResult{
		Validated: true,
		Message:   "",
	}

	vr.SetFailed("lorem ipsum")

	if vr.Validated {
		t.Errorf("validation result should be false")
	}

	if vr.Message != "lorem ipsum, " {
		t.Errorf("invalid validation message: %s", vr.Message)
	}
}
