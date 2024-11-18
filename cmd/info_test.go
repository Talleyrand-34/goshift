package cmd_test

import (
	"encoding/json"
	"goshift/cmd"
	"testing"

	"github.com/containerd/btrfs/v2"
	"github.com/stretchr/testify/assert"
)

func TestBasicInfo(t *testing.T) {
	testSubvolPath := SubvolPath1
	BasicTestInfo(t, testSubvolPath)
}

func TestRoot(t *testing.T) {
	testSubvolPath := RootPath
	BasicTestInfo(t, testSubvolPath)
}
func TestNotAPath(t *testing.T) {
	testSubvolPath := "/notasubvolume"
	cmdjson, err := cmd.BasicInfo([]string{testSubvolPath})
	assert.Error(t, err, "Error: no such file or directory")
	assert.Nil(t, cmdjson, "Expected nil output for non-subvolume path")
}
func TestNotASubvolumeButSubfolder(t *testing.T) {
	testSubvolPath := SubFolderPath
	cmdjson, err := cmd.BasicInfo([]string{testSubvolPath})
	assert.EqualError(t, err, testSubvolPath+" is a btrfs subfolder", "Expected specific error message for non-subvolume path")
	assert.Nil(t, cmdjson, "Expected nil output for non-subvolume path")
}
func TestNotBtrfs(t *testing.T) {
	testSubvolPath := MountPath
	cmdjson, err := cmd.BasicInfo([]string{testSubvolPath})
	assert.EqualError(t, err, testSubvolPath+" is not a valid btrfs subvolume", "Expected specific error message for non-subvolume path")
	assert.Nil(t, cmdjson, "Expected nil output for non-subvolume path")
}

func BasicTestInfo(t *testing.T, testSubvolPath string) {
	// Call BasicInfo with our test subvolume
	cmdjson, err := cmd.BasicInfo([]string{testSubvolPath})
	assert.NoError(t, err, "Failed to get subvolume info fn")

	// Get the info directly using btrfs package to compare
	info, err := btrfs.SubvolInfo(testSubvolPath)
	assert.NoError(t, err, "Failed to get subvolume info manually")

	// Verify the info structure has expected fields
	assert.NotZero(t, info.ID, "Subvolume ID should not be zero")
	assert.NotEmpty(t, info.UUID, "UUID should not be empty")

	// Verify the info can be marshaled to JSON
	manualjson, err := json.Marshal(info)
	assert.NoError(t, err, "Failed to marshal subvolume info to JSON")
	assert.Equal(t, string(manualjson), string(cmdjson))
}
