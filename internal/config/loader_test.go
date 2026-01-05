package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func setupTests(t *testing.T) {
	t.Helper()
	_ = os.Unsetenv("DB_NAME")
	_ = os.Unsetenv("DB_USERNAME")
	_ = os.Unsetenv("DB_PASSWORD")
	_ = os.Unsetenv("DB_HOST")
	_ = os.Unsetenv("DB_PORT")
	_ = os.Unsetenv("JWT_SECRET")
}

type mockSetup struct {
	prepareDirFn func(t *testing.T)
}

func (m mockSetup) prepareDir(t *testing.T) {
	t.Helper()
	if m.prepareDirFn != nil {
		m.prepareDirFn(t)
	}
}

var validEnv = `DB_NAME=test
DB_USERNAME=testingUser
DB_PASSWORD=secretPassword
DB_HOST="localhost"
DB_PORT=5432

JWT_SECRET="secret-key-for-testing"`

func TestLoad(t *testing.T) {
	var tests = []struct {
		name        string
		mock        mockSetup
		expectedCfg Config
		shouldFail  bool
	}{
		{
			"valid env file",
			mockSetup{prepareDirFn: func(t *testing.T) {
				tmpDir := t.TempDir()
				err := os.WriteFile(filepath.Join(tmpDir, ".env"), []byte(validEnv), 0644)
				if err != nil {
					t.Fatal(err)
				}
				orgDir, _ := os.Getwd()
				_ = os.Chdir(tmpDir)
				t.Cleanup(func() {
					_ = os.Chdir(orgDir)
				})
			}},
			Config{
				DB: DBConfig{
					Name:     "test",
					User:     "testingUser",
					Password: "secretPassword",
					Host:     "localhost",
					Port:     "5432",
				},
				JWT: JWTConfig{
					Secret: "secret-key-for-testing",
				},
			},
			false,
		},
		{
			"additional value in env file",
			mockSetup{prepareDirFn: func(t *testing.T) {
				var invalidEnv = `
				DB_NAME="test"
				DB_USERNAME=testingUser
				DB_PASSWORD=secretPassword
				DB_HOST="localhost"
				DB_PORT=5432
				SOME_VAL=123
				
				JWT_SECRET="secret-key-for-testing"`
				tmpDir := t.TempDir()
				err := os.WriteFile(filepath.Join(tmpDir, ".env"), []byte(invalidEnv), 0644)
				if err != nil {
					t.Fatal(err)
				}
				orgDir, _ := os.Getwd()
				_ = os.Chdir(tmpDir)
				t.Cleanup(func() {
					_ = os.Chdir(orgDir)
				})
			}},
			Config{
				DB: DBConfig{
					Name:     "test",
					User:     "testingUser",
					Password: "secretPassword",
					Host:     "localhost",
					Port:     "5432",
				},
				JWT: JWTConfig{
					Secret: "secret-key-for-testing",
				},
			},
			false,
		},
		{
			"missing values from env file",
			mockSetup{prepareDirFn: func(t *testing.T) {
				tmpDir := t.TempDir()
				err := os.WriteFile(filepath.Join(tmpDir, ".env"), []byte(""), 0644)
				if err != nil {
					t.Fatal(err)
				}
				orgDir, _ := os.Getwd()
				_ = os.Chdir(tmpDir)
				t.Cleanup(func() {
					_ = os.Chdir(orgDir)
				})
			}},
			Config{},
			true,
		},
		{
			"missing single value from env file",
			mockSetup{prepareDirFn: func(t *testing.T) {
				var invalidEnv = `
				DB_USERNAME=testingUser
				DB_PASSWORD=secretPassword
				DB_HOST="localhost"
				DB_PORT=5432
				
				JWT_SECRET="secret-key-for-testing"`
				tmpDir := t.TempDir()
				err := os.WriteFile(filepath.Join(tmpDir, ".env"), []byte(invalidEnv), 0644)
				if err != nil {
					t.Fatal(err)
				}
				orgDir, _ := os.Getwd()
				_ = os.Chdir(tmpDir)
				t.Cleanup(func() {
					_ = os.Chdir(orgDir)
				})
			}},
			Config{},
			true,
		},
		{
			"missing jwt secret value from env file",
			mockSetup{prepareDirFn: func(t *testing.T) {
				var invalidEnv = `
				DB_NAME=test
				DB_USERNAME=testingUser
				DB_PASSWORD=secretPassword
				DB_HOST="localhost"
				DB_PORT=5432`
				tmpDir := t.TempDir()
				err := os.WriteFile(filepath.Join(tmpDir, ".env"), []byte(invalidEnv), 0644)
				if err != nil {
					t.Fatal(err)
				}
				orgDir, _ := os.Getwd()
				_ = os.Chdir(tmpDir)
				t.Cleanup(func() {
					_ = os.Chdir(orgDir)
				})
			}},
			Config{},
			true,
		},
		{
			"incorrect value in env file",
			mockSetup{prepareDirFn: func(t *testing.T) {
				var invalidEnv = `
				DB_NAME="   "
				DB_USERNAME=testingUser
				DB_PASSWORD=secretPassword
				DB_HOST="localhost"
				DB_PORT=5432
				
				JWT_SECRET="secret-key-for-testing"`
				tmpDir := t.TempDir()
				err := os.WriteFile(filepath.Join(tmpDir, ".env"), []byte(invalidEnv), 0644)
				if err != nil {
					t.Fatal(err)
				}
				orgDir, _ := os.Getwd()
				_ = os.Chdir(tmpDir)
				t.Cleanup(func() {
					_ = os.Chdir(orgDir)
				})
			}},
			Config{},
			true,
		},
		{
			"missing env file entirely",
			mockSetup{prepareDirFn: func(t *testing.T) {
			}},
			Config{},
			true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			setupTests(t)
			tc.mock.prepareDir(t)

			failed := false
			fatalf = func(format string, args ...any) {
				failed = true
				panic(fmt.Sprintf(format, args...))
			}
			t.Cleanup(func() {
				fatalf = log.Fatalf
			})

			defer func() {
				recover()
				if tc.shouldFail && !failed {
					t.Errorf("expected Load to fail but it did not")
				}
				if !tc.shouldFail && failed {
					t.Errorf("Load failed unexpectedly")
				}
			}()

			cfg := Load()

			if diff := cmp.Diff(tc.expectedCfg, cfg); diff != "" {
				t.Errorf("unexpected config loaded <-want,+got>\n%s", diff)
			}
		})
	}
}
