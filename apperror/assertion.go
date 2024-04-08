package apperror

import "testing"

func AssertErrorIsCode(t *testing.T, err error, code int) {
	if !IsErrorCode(err, code) {
		t.Errorf("AssertErrorIsCode failed: want code: %d, got: %v", code, err)
	}
}
