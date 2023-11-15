package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsZero(t *testing.T) {
	type mockStruct struct {
		A string
		B int
	}
	cases := []Case{
		{"nil", nil, true},
		{"arr", []interface{}{1, "2"}, false},
		{"empty struct", mockStruct{}, true},
		{"int 0", 0, true},
		{"int", 1, false},
	}
	t.Parallel()
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			assert.Equal(t, tc.Input, tc.Expect)
		})
	}
}

type Case struct {
	Name   string
	Input  interface{}
	Expect interface{}
}
