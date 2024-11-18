package cmd_test

import (
	"goshift/cmd"
	"os"
	"testing"

	"github.com/containerd/btrfs/v2"
	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/assert"
)

func TestGetMountpoint(t *testing.T) {

	tests := []struct {
		subvolPath string
		phdev      string
		subv       string
		err        error
	}{
		{SubvolPath1, "/dev/loop0", "[/@testsubvol]", nil},
	}

	for _, test := range tests {
		phDev, subv, err := cmd.GetMountpoint(test.subvolPath)
		assert.Equal(t, test.phdev, phDev)
		assert.Equal(t, test.subv, subv)
		assert.Equal(t, test.err, err)
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

func TestCreateSubvolumeRedhatStyle(t *testing.T) {
	// Test cases where
	// - RootPath + @ + subvolname
	// - MountPath + subvolname
	SubRootPath1 := RootPath + "/@rh1"
	SubMountPath1 := MountPath + "/rh1"
	SubRootPath2 := RootPath + "/@rh2"
	SubMountPath2 := MountPath + "/rh2"
	SubRootPath3 := RootPath + "/@rh3"
	SubMountPath3 := MountPath + "/rh3"
	tests := []struct {
		name          string
		args          string
		createErr     error
		expectedError bool
		expectedPath  string
	}{
		{
			name:          "Normal create success",
			args:          SubMountPath1,
			createErr:     nil,
			expectedError: false,
			expectedPath:  SubRootPath1,
		}, {
			name:          "Normal create success",
			args:          SubMountPath2,
			createErr:     nil,
			expectedError: false,
			expectedPath:  SubRootPath2,
		}, {
			name:          "Normal create success",
			args:          SubMountPath3,
			createErr:     nil,
			expectedError: false,
			expectedPath:  SubRootPath3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := cmd.CreateSubvolumeRedhatStyle(test.args)
			assert.Equal(t, test.expectedError, err != nil)
			//ensure path is correct
			_, err = os.Stat(test.expectedPath)
			assert.NoError(t, err, "Expected folder path to exist: %s", test.expectedPath)
			err = btrfs.IsSubvolume(SubRootPath1)
			assert.NoError(t, err, "Expected subvolume path to exist: %s", test.expectedPath)

		})
	}
}

/* func TestInterfaceCreateSubvolume(t *testing.T) {
	testCases := []struct {
		name          string
		args          []string
		createErr     error
		redhatCalled  bool
		expectedError bool
	}{
		{
			name:          "Normal create success",
			args:          []string{SubvolPath1},
			createErr:     nil,
			redhatCalled:  false,
			expectedError: false,
		},
		{
			name:          "Create error",
			args:          []string{MountPath},
			createErr:     assert.AnError,
			redhatCalled:  false,
			expectedError: true,
		},
		{
			name:          "Redhat style",
			args:          []string{SubvolPath1 + "[redhat]"},
			createErr:     nil,
			redhatCalled:  true,
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			redhatCalled := false
			mockCreate := func(path string) error {
				return tc.createErr
			}
			mockRedhat := func(path string) {
				redhatCalled = true
			}

			// Note: These functions need to be exported in the cmd package to be testable
			interface_create_subvolume(nil, tc.args, mockCreate, mockRedhat)

			assert.Equal(t, tc.redhatCalled, redhatCalled, "Redhat style creation mismatch")
		})
	}
}

func TestCreateSubvolumeRedhatStyle(t *testing.T) {
	// Since create_subvolume_redhat_style doesn't return anything and just calls other functions,
	// we can only verify it doesn't panic
	assert.NotPanics(t, func() {
		create_subvolume_redhat_style(SubvolPath1)
	})
} */
