package models

import (
	"math"
	"strconv"
	"testing"
)

func TestPriority_UnmarshalJSON(t *testing.T) {
	var tests = []struct {
		name           string
		testedVal      string
		wantsError     bool
		expectedError  string
		wantedPriority Priority
	}{
		{
			"string value low",
			`"low"`,
			false,
			"",
			PriorityLow,
		},
		{
			"string value medium",
			`"medium"`,
			false,
			"",
			PriorityMedium,
		},
		{
			"string value high",
			`"high"`,
			false,
			"",
			PriorityHigh,
		},
		{
			"number value low",
			"0",
			false,
			"",
			PriorityLow,
		},
		{
			"number value medium",
			"1",
			false,
			"",
			PriorityMedium,
		},
		{
			"number value high",
			"2",
			false,
			"",
			PriorityHigh,
		},
		{
			"incorrect string value",
			`"lorem"`,
			true,
			"invalid priority value: \"lorem\"",
			0,
		},
		{
			"incorrect number value",
			"12",
			true,
			"invalid priority value: 12",
			0,
		},
		{
			"entirely wrong string",
			"lorem",
			true,
			"invalid payload: [108 111 114 101 109]",
			0,
		},
		{
			"empty string",
			"",
			true,
			"invalid payload: []",
			0,
		},
		{
			"whitespaces",
			"   ",
			true,
			"invalid payload: [32 32 32]",
			0,
		},
		{
			"negative number",
			"-2",
			true,
			"invalid priority value: -2",
			0,
		},
		{
			"decimal number",
			"1.1",
			true,
			"invalid payload: [49 46 49]",
			0,
		},
		{
			"very large number",
			strconv.Itoa(math.MaxInt),
			true,
			"invalid priority value: " + strconv.Itoa(math.MaxInt),
			0,
		},
	}

	for _, i := range tests {
		t.Run(i.name, func(t *testing.T) {
			var p Priority
			err := p.UnmarshalJSON([]byte(i.testedVal))

			if i.wantsError && err == nil {
				t.Errorf("for value %s it should return an error", i.testedVal)
			}

			if !i.wantsError && err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if i.wantsError {
				if err.Error() != i.expectedError {
					t.Errorf("expected error <%s> but got <%s>", i.expectedError, err)
				}
				if p != 0 {
					t.Errorf("UnmarshalJSON returned an error but priority value changed")
				}
			} else {
				if p != i.wantedPriority {
					t.Errorf("expected value %d but got %d", i.wantedPriority, p)
				}
			}
		})
	}
}

func TestPriority_UnmarshalJSON_nil_value(t *testing.T) {
	var p Priority
	expectedError := "invalid payload: []"
	err := p.UnmarshalJSON([]byte(nil))
	if err == nil {
		t.Errorf("function should return an error but it did not")
	} else {
		if err.Error() != expectedError {
			t.Errorf("expected error <%s> but got <%s>", expectedError, err)
		}
	}

}
