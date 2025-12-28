package helpers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func testSliceContainsGeneric[T comparable](t *testing.T, tests []struct {
	name           string
	expectedResult bool
	val            T
	slice          []T
}) {
	t.Helper()
	for _, i := range tests {
		t.Run(i.name, func(t *testing.T) {
			exists := SliceContains(i.slice, i.val)
			if exists != i.expectedResult {
				t.Errorf("expected result to be %t but got %t instead", i.expectedResult, exists)
			}
		})
	}
}

func TestSliceContainsString(t *testing.T) {
	tests := []struct {
		name           string
		expectedResult bool
		val            string
		slice          []string
	}{
		{
			"slice contains",
			true,
			"lorem",
			[]string{"lorem", "ipsum"},
		},
		{
			"slice does not contain",
			false,
			"dolor",
			[]string{"lorem", "ipsum"},
		},
		{
			"empty slice",
			false,
			"lorem",
			[]string{},
		},
		{
			"empty search value",
			false,
			"",
			[]string{"lorem", "ipsum", "dolor"},
		},
		{
			"duplicate values in slice",
			true,
			"lorem",
			[]string{"lorem", "lorem", "ipsum", "dolor"},
		},
		{
			"empty string in slice",
			true,
			"lorem",
			[]string{"lorem", "", "ipsum", "dolor"},
		},
		{
			"empty string in val and slice",
			true,
			"",
			[]string{"", "lorem", "ipsum"},
		},
	}
	testSliceContainsGeneric(t, tests)
}

func TestSliceContainsBool(t *testing.T) {
	var tests = []struct {
		name           string
		expectedResult bool
		val            bool
		slice          []bool
	}{
		{
			"slice contains",
			true,
			true,
			[]bool{true, false},
		},
		{
			"slice does not contain",
			false,
			true,
			[]bool{false},
		},
		{
			"duplicate value",
			true,
			true,
			[]bool{true, false, true},
		},
		{
			"empty slice",
			false,
			true,
			[]bool{},
		},
	}

	testSliceContainsGeneric(t, tests)
}

func TestSliceContainsInt(t *testing.T) {
	tests := []struct {
		name           string
		expectedResult bool
		val            int
		slice          []int
	}{
		{
			"slice contains",
			true,
			12,
			[]int{1, 12, 123},
		},
		{
			"slice does not contain",
			false,
			21,
			[]int{1, 2, 3, 4},
		},
		{
			"empty slice",
			false,
			1,
			[]int{},
		},
		{
			"duplicate values in slice",
			true,
			21,
			[]int{21, 37, 21, 11},
		},
	}
	testSliceContainsGeneric(t, tests)
}

func TestSliceContainsPoints(t *testing.T) {
	type point struct {
		x, y int
	}
	var tests = []struct {
		name           string
		expectedResult bool
		val            point
		slice          []point
	}{
		{
			"slice contains",
			true,
			point{1, 2},
			[]point{
				{1, 2},
				{3, 4},
				{5, 6},
			},
		},
		{
			"slice does not contain",
			false,
			point{1, 2},
			[]point{
				{3, 4},
				{5, 6},
			},
		},
		{
			"duplicate value in slice",
			true,
			point{1, 2},
			[]point{
				{1, 2},
				{1, 2},
				{3, 4},
			},
		},
		{
			"empty search value",
			false,
			point{},
			[]point{
				{21, 37},
			},
		},
	}

	testSliceContainsGeneric(t, tests)
}

func TestSliceContainsCustomType(t *testing.T) {
	type UserId int
	var tests = []struct {
		name           string
		expectedResult bool
		val            UserId
		slice          []UserId
	}{
		{
			"slice contains",
			true,
			1,
			[]UserId{1, 2, 3},
		},
		{
			"slice does not contain",
			false,
			12,
			[]UserId{1, 2, 3},
		},
		{
			"empty slice",
			false,
			1,
			[]UserId{},
		},
	}

	testSliceContainsGeneric(t, tests)
}

func TestSliceContainsConst(t *testing.T) {
	type status int
	const (
		statusQueued status = iota
		statusRunning
		statusDone
		statusFailed
	)
	var userStatus status = 123

	var tests = []struct {
		name           string
		expectedResult bool
		val            status
		slice          []status
	}{
		{
			"slice contains",
			true,
			statusDone,
			[]status{statusQueued, statusRunning, statusDone},
		},
		{
			"slice does not contain",
			false,
			statusQueued,
			[]status{statusRunning, statusDone, statusFailed},
		},
		{
			"search value out of range",
			false,
			userStatus,
			[]status{statusQueued, statusRunning, statusRunning},
		},
	}

	testSliceContainsGeneric(t, tests)
}

func TestSliceContainsPointers(t *testing.T) {
	v1 := 21
	v2 := 37
	v3 := 21

	ptr1 := &v1
	ptr2 := &v2
	ptr3 := &v3
	ptrCpy := ptr1

	var tests = []struct {
		name           string
		expectedResult bool
		val            *int
		slice          []*int
	}{
		{
			"slice contains",
			true,
			ptr1,
			[]*int{ptr1, ptr2},
		},
		{
			"slice does not contain",
			false,
			ptr2,
			[]*int{ptr1},
		},
		{
			"slice contains via copy",
			true,
			ptrCpy,
			[]*int{ptr1, ptr2},
		},
		{
			"slice does not contains - same val different address",
			false,
			ptr3,
			[]*int{ptr1, ptr2},
		},
		{
			"empty slice",
			false,
			ptr1,
			[]*int{},
		},
		{
			"nil pointer not in slice",
			false,
			nil,
			[]*int{ptr1, ptr2},
		},
		{
			"nil pointer in slice",
			true,
			nil,
			[]*int{nil, ptr2},
		},
	}

	testSliceContainsGeneric(t, tests)
}

func TestSliceContainsArrays(t *testing.T) {
	type pointsArray [2]int
	var tests = []struct {
		name           string
		expectedResult bool
		val            pointsArray
		slice          []pointsArray
	}{
		{
			"slice contains",
			true,
			pointsArray{1, 2},
			[]pointsArray{
				{1, 2},
				{3, 4},
			},
		},
		{
			"slice does not contains",
			false,
			pointsArray{1, 2},
			[]pointsArray{
				{3, 4},
				{21, 37},
			},
		},
		{
			"empty slice",
			false,
			pointsArray{1, 2},
			[]pointsArray{},
		},
		{
			"empty val",
			false,
			pointsArray{},
			[]pointsArray{
				{1, 2},
				{3, 4},
			},
		},
		{
			"part match",
			false,
			pointsArray{1, 4},
			[]pointsArray{
				{1, 2},
				{3, 4},
			},
		},
		{
			"duplicates in slice",
			true,
			pointsArray{1, 2},
			[]pointsArray{
				{1, 2},
				{1, 2},
				{3, 4},
			},
		},
	}

	testSliceContainsGeneric(t, tests)
}

func TestGetQueryPath(t *testing.T) {
	var tests = []struct {
		name        string
		passedStr   string
		expectedStr string
	}{
		{
			"simple file name",
			"lorem.ips",
			"./internal/db/queries/lorem.ips",
		},
		{
			"simple word",
			"lorem",
			"./internal/db/queries/lorem",
		},
		{
			"empty string",
			"",
			"./internal/db/queries/",
		},
		{
			"special chars",
			"l0r@m.#$@",
			"./internal/db/queries/l0r@m.#$@",
		},
		{
			"path separator",
			"//lorem.ips",
			"./internal/db/queries///lorem.ips",
		},
	}

	for _, i := range tests {
		t.Run(i.name, func(t *testing.T) {
			qPath := GetQueryPath(i.passedStr)
			if qPath != i.expectedStr {
				t.Errorf("returned path <%s> but expected value <%s>", qPath, i.expectedStr)
			}
		})
	}
}

func TestHashPassword(t *testing.T) {
	var tests = []struct {
		name         string
		password     string
		expectsError bool
	}{
		{
			"simple password",
			"loremipsum",
			false,
		},
		{
			"short password",
			"a",
			false,
		},
		{
			"long password",
			strings.Repeat("a", 73),
			true,
		},
		{
			"empty string",
			"",
			false,
		},
	}

	for _, i := range tests {
		t.Run(i.name, func(t *testing.T) {
			h, err := HashPassword(i.password)
			if i.expectsError && err == nil {
				t.Errorf("expected error while hashing but got nil")
			}

			if !i.expectsError && err != nil {
				t.Errorf("unexpected error while hashing password: %s", err)
			}

			if h == i.password || strings.EqualFold(h, i.password) {
				t.Errorf("hashed password can't be the same as original value")
			}

			if !i.expectsError {
				err = bcrypt.CompareHashAndPassword([]byte(h), []byte(i.password))
				if err != nil {
					t.Errorf("unexpected error while comparing values: %s", err)
				}
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	validPassword := "loremipsum"
	validHash, _ := HashPassword(validPassword)

	diffPass := "doloret"
	diffHash, _ := HashPassword(diffPass)

	emptyPassword := ""
	emptyHash, _ := HashPassword(emptyPassword)

	nonEmptyHash, _ := HashPassword("nonempty")

	var tests = []struct {
		name           string
		password       string
		hash           string
		expectedResult bool
	}{
		{
			"valid password and hash",
			validPassword,
			validHash,
			true,
		},
		{
			"wrong password, valid hash",
			diffPass,
			validHash,
			false,
		},
		{
			"correct password, wrong hash",
			validPassword,
			diffHash,
			false,
		},
		{
			"empty password, empty hash",
			emptyPassword,
			emptyHash,
			true,
		},
		{
			"empty password, not empty hash",
			emptyPassword,
			nonEmptyHash,
			false,
		},
		{
			"not empty password, empty hash",
			validPassword,
			emptyHash,
			false,
		},
		{
			"valid password, malformed hash",
			validPassword,
			"$2a$06$invalidhashformat",
			false,
		},
		{
			"valid password, too short hash",
			validPassword,
			"$2a$06$short",
			false,
		},
		{
			"valid password, completely wrong hash",
			validPassword,
			"totalywronghash",
			false,
		},
	}

	for _, i := range tests {
		t.Run(i.name, func(t *testing.T) {
			r := ValidatePassword(i.password, i.hash)

			if r != i.expectedResult {
				t.Errorf("expected result %t but got %t", i.expectedResult, r)
			}
		})
	}
}

func TestJsonResponse(t *testing.T) {
	var tests = []struct {
		name           string
		status         int
		payload        interface{}
		expectedStatus int
		expectBody     bool
		expectedBody   string
	}{
		{
			"success request",
			http.StatusOK,
			nil,
			http.StatusOK,
			false,
			"",
		},
		{
			"failed with message",
			http.StatusUnprocessableEntity,
			map[string]string{"error": "email is required"},
			http.StatusUnprocessableEntity,
			true,
			`{"error":"email is required"}`,
		},
		{
			"empty body",
			http.StatusNoContent,
			nil,
			http.StatusNoContent,
			false,
			"",
		},
		{
			"success with payload",
			http.StatusOK,
			map[string]string{
				"message": "ok",
			},
			http.StatusOK,
			true,
			`{"message":"ok"}`,
		},
		{
			"empty map payload",
			http.StatusOK,
			map[string]string{},
			http.StatusOK,
			true,
			"{}",
		},
		{
			"string in payload",
			http.StatusOK,
			"lorem",
			http.StatusOK,
			true,
			"\"lorem\"",
		},
		{
			"struct in payload",
			http.StatusOK,
			struct {
				name string
				code int
			}{"lorem", 12},
			http.StatusOK,
			true,
			"{}",
		},
		{
			"slice in payload",
			http.StatusOK,
			[]int{1, 2, 3},
			http.StatusOK,
			true,
			"[1,2,3]",
		},
	}

	for _, i := range tests {
		t.Run(i.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			JsonResponse(w, i.status, i.payload)

			ct := w.Header().Get("Content-Type")
			if ct != "application/json" {
				t.Errorf("expected content type to be application/json but got %s", ct)
			}

			if w.Code != i.expectedStatus {
				t.Errorf("expected status code %d but got %d", i.expectedStatus, w.Code)
			}

			b := w.Body.String()

			if i.expectBody {
				b = strings.TrimSpace(b)
				if b != i.expectedBody {
					t.Errorf("body mismatch, expected %s but got %s", i.expectedBody, b)
				}
			} else {
				if b != "" {
					t.Errorf("expected empty body but got %s", b)
				}
			}
		})
	}
}

func TestIsValidEmail(t *testing.T) {
	var tests = []struct {
		name           string
		email          string
		expectedResult bool
	}{
		{
			"valid email",
			"lorem@ipsum.com",
			true,
		},
		{
			"invalid string",
			"lorem ipsum",
			false,
		},
		{
			"plus sign in address",
			"lorem+ipsum@dolor.et",
			true,
		},
		{
			"email with special char",
			"lorem.ipsum+dolor@est.co.uk",
			true,
		},
		{
			"email with subdomain",
			"lorem@ipsum.dolor.com",
			true,
		},
		{
			"invalid format - only domain",
			"@ipsum.dolor.com",
			false,
		},
		{
			"invalid format - only username",
			"lorem@",
			false,
		},
		{
			"invalid format - missing domain",
			"lorem@com",
			false,
		},
		{
			"invalid format - missing domain with dot",
			"lorem@.com",
			false,
		},
		{
			"invalid format - missing country",
			"lorem@ipsum",
			false,
		},
		{
			"email with spaces",
			"lorem ipsum@dolor",
			false,
		},
		{
			"email with quotes",
			"lorem\"ipsum\"@dolor.com",
			false,
		},
	}

	for _, i := range tests {
		t.Run(i.name, func(t *testing.T) {
			ok := IsValidEmail(i.email)
			if ok != i.expectedResult {
				t.Errorf("%s: failed, function returned %t but %t was expected", i.name, ok, i.expectedResult)
			}
		})
	}
}
