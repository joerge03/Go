package main

import (
	"testing"
)

func TestIsValidEmail(t *testing.T) {
	testData := []struct {
		email string
		want  bool
	}{
		{"testing@gmail.com", false},
		{"testing@sdfasdf", false},
	}

	for _, data := range testData {
		result := IsValidEmail(data.email)
		if data.want != result {
			t.Errorf("IsValidEmail(%v)= %v , expected = %v", data.email, result, data.want)
		}
	}
}
