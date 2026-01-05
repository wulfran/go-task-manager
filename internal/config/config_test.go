package config

import "testing"

func TestDBConfig_Validate(t *testing.T) {
	var tests = []struct {
		name         string
		dbStruct     DBConfig
		expectsError bool
		errorWanted  string
	}{
		{
			"valid struct, no errors",
			DBConfig{
				Name:     "Lorem",
				User:     "Ipsum",
				Password: "Dolor",
				Host:     "Et",
				Port:     "Amet",
			},
			false,
			"",
		},
		{
			"name missing",
			DBConfig{
				Name:     "",
				User:     "Ipsum",
				Password: "Dolor",
				Host:     "Et",
				Port:     "Amet",
			},
			true,
			"Name is required",
		},
		{
			"user missing",
			DBConfig{
				Name:     "Lorem",
				User:     "",
				Password: "Dolor",
				Host:     "Et",
				Port:     "Amet",
			},
			true,
			"User is required",
		},
		{
			"Password missing",
			DBConfig{
				Name:     "Lorem",
				User:     "Ipsum",
				Password: "",
				Host:     "Et",
				Port:     "Amet",
			},
			true,
			"Password is required",
		},
		{
			"Host missing",
			DBConfig{
				Name:     "Lorem",
				User:     "Ipsum",
				Password: "Dolor",
				Host:     "",
				Port:     "Amet",
			},
			true,
			"Host is required",
		},
		{
			"Port missing",
			DBConfig{
				Name:     "Lorem",
				User:     "Ipsum",
				Password: "Dolor",
				Host:     "Et",
				Port:     "",
			},
			true,
			"Port is required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.dbStruct.Validate()

			if tc.expectsError && err == nil {
				t.Errorf("validate should return an error but it did not")
			}
			if !tc.expectsError && err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if tc.expectsError && err != nil {
				if tc.errorWanted != err.Error() || err.Error() == "" {
					t.Errorf("wrong error returned, expected <%s> but got <%s>", tc.errorWanted, err)
				}
			}
		})
	}
}

func TestJWTConfig_Validate(t *testing.T) {
	var tests = []struct {
		name         string
		jwtStruct    JWTConfig
		expectsError bool
		errorWanted  string
	}{
		{
			"Valid input, no errors",
			JWTConfig{Secret: "test-secret"},
			false,
			"",
		},
		{
			"Secret missing",
			JWTConfig{Secret: ""},
			true,
			"Secret is required",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.jwtStruct.Validate()

			if tc.expectsError && err == nil {
				t.Errorf("validate should return an error but it did not")
			}
			if !tc.expectsError && err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if tc.expectsError && err != nil {
				if tc.errorWanted != err.Error() || err.Error() == "" {
					t.Errorf("wrong error returned, expected <%s> but got <%s>", tc.errorWanted, err)
				}
			}
		})
	}
}
