package cmd_test

import (
	"goshift/cmd"
	"testing"

	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/assert"
)

func TestGetMountpoint(t *testing.T) {
	// Create test subvolume path
	//testMountPath := MountPath
	testSubvolPath := SubvolPath1

	// Test GetMountpoint function
	phDev, subv, err := cmd.GetMountpoint(testSubvolPath)
	if err != nil {
		t.Fatalf("GetMountpoint failed: %v", err)
	}

	// Verify physical device is not empty
	if phDev == "" {
		t.Error("Expected physical device path, got empty string")
	}

	// Verify subvolume info is returned
	if subv == "" {
		t.Error("Expected subvolume info, got empty string")
	}

}

//func TestSplitStringWithBrackets(t *testing.T) {
//	input := "mainstring[substring]"
//	expectedMain := "mainstring"
//	expectedSub := "[substring]"
//	main, sub := cmd.SplitStringWithBrackets(input)
//	if main != expectedMain || sub != expectedSub {
//		t.Errorf("SplitStringWithBrackets() = (%v, %v), want (%v, %v)", main, sub, expectedMain, expectedSub)
//	}
//}

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
