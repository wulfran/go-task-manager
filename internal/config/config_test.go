package config

import "testing"

func testValidateStruct[Struct testStruct](t *testing.T, s Struct, expectsError bool, errorWanted string) {
	t.Helper()
	err := s.Validate()

	if expectsError && err == nil {
		t.Errorf("validate should return an error but it did not")
	}
	if !expectsError && err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if expectsError && err != nil {
		if errorWanted != err.Error() {
			t.Errorf("wrong error returned, expected <%s> but got <%s>", errorWanted, err)
		}
	}
}

type testStruct interface {
	DBConfig |
		JWTConfig |
		structWithInt
	Validate() error
}

type structWithInt struct {
	name  string
	count int
}

func (s structWithInt) Validate() error {
	return validateStruct(s)
}

func TestDBConfig_Validate(t *testing.T) {
	t.Parallel()
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
				Password: "Dolor",
				Host:     "Et",
				Port:     "Amet",
			},
			true,
			"User is required",
		},
		{
			"password missing",
			DBConfig{
				Name: "Lorem",
				User: "Ipsum",
				Host: "Et",
				Port: "Amet",
			},
			true,
			"Password is required",
		},
		{
			"host missing",
			DBConfig{
				Name:     "Lorem",
				User:     "Ipsum",
				Password: "Dolor",
				Port:     "Amet",
			},
			true,
			"Host is required",
		},
		{
			"port missing",
			DBConfig{
				Name:     "Lorem",
				User:     "Ipsum",
				Password: "Dolor",
				Host:     "Et",
			},
			true,
			"Port is required",
		},
		{
			"name is a whitespace",
			DBConfig{
				Name:     "   ",
				User:     "Ipsum",
				Password: "Dolor",
				Host:     "Et",
				Port:     "Amet",
			},
			true,
			"Name is required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			testValidateStruct(t, tc.dbStruct, tc.expectsError, tc.errorWanted)
		})
	}
}

func TestJWTConfig_Validate(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		name         string
		jwtStruct    JWTConfig
		expectsError bool
		errorWanted  string
	}{
		{
			"valid input, no errors",
			JWTConfig{Secret: "test-secret"},
			false,
			"",
		},
		{
			"secret missing",
			JWTConfig{},
			true,
			"Secret is required",
		},
		{
			"secret is whitespace",
			JWTConfig{Secret: "  "},
			true,
			"Secret is required",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			testValidateStruct(t, tc.jwtStruct, tc.expectsError, tc.errorWanted)
		})
	}
}

func TestValidateZeroValue(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		name         string
		inputStruct  structWithInt
		expectsError bool
		errorWanted  string
	}{
		{
			"valid value",
			structWithInt{
				name:  "Lorem",
				count: 12,
			},
			false,
			"",
		},
		{
			"zero value",
			structWithInt{
				name:  "Lorem",
				count: 0,
			},
			true,
			"count is required",
		},
		{
			"missing int value",
			structWithInt{
				name: "Lorem",
			},
			true,
			"count is required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			testValidateStruct(t, tc.inputStruct, tc.expectsError, tc.errorWanted)
		})
	}
}
