package cmd_test

import (
	"goshift/cmd"
	"testing"

	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/assert"
)

func TestSplitStringWithBrackets(t *testing.T) {
	tests := []struct {
		input        string
		expectedMain string
		expectedSub  string
	}{
		{"device[0]", "device", "[0]"},
		{"device", "device", ""},
	}

	for _, test := range tests {
		// mainStr, subStr := SplitStringWithBrackets(test.input)
		mainStr, subStr := cmd.SplitStringWithBrackets(test.input)
		assert.Equal(t, test.expectedMain, mainStr)
		assert.Equal(t, test.expectedSub, subStr)
	}
}
