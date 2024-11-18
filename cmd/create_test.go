package cmd_test

import (
	"fmt"
	"goshift/cmd"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/containerd/btrfs/v2"
	"github.com/kgs19/cmdx"

	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/assert"
)

func createBtrfsImage() string {
	// Create a temporary file to use as a Btrfs filesystem
	tempFile, err := os.CreateTemp("", "btrfs_test.img")
	if err != nil {
		fmt.Printf("Error creating temporary file: %v\n", err)
		os.Exit(1)
	}

	// Set the size of the file (e.g., 100MB)
	if err := tempFile.Truncate(1114294784); err != nil {
		fmt.Printf("Error setting size of temporary file: %v\n", err)
		os.Exit(1)
	}
	tempFile.Close()

	// Format the file as Btrfs
	if err := cmdx.RunCommandPrintOutput("mkfs.btrfs", tempFile.Name()); err != nil {
		fmt.Printf("Error formatting file as Btrfs: %v\n", err)
		os.Exit(1)
	}

	return tempFile.Name()
}

func createTestBtrfsMount(imagePath string) string {
	// Create a temporary mount directory
	testMountPath := filepath.Join(os.TempDir(), "testbtrfs")
	if err := os.MkdirAll(testMountPath, 0744); err != nil {
		fmt.Printf("Error creating temporary mount directory: %v\n", err)
		os.Exit(1)
	}

	// Mount the temporary Btrfs filesystem
	if err := cmdx.RunCommandPrintOutput("mount", imagePath, testMountPath); err != nil {
		fmt.Printf("Error mounting Btrfs filesystem: %v\n", err)
		os.Exit(1)
	}

	return testMountPath
}

func createTestSubvolume(mountPath string, name string) string {
	testSubvolPath := filepath.Join(mountPath, name)

	// Create test subvolume
	if err := btrfs.SubvolCreate(testSubvolPath); err != nil {
		fmt.Printf("Failed to create test subvolume: %v", err)
		os.Exit(1)
	}

	return testSubvolPath
}

func createFalseMountpoint(mountPathOrigin string, mountPathDest string) string {
	if err := os.MkdirAll(mountPathDest, 0744); err != nil {
		fmt.Printf("Error creating mountPathOrigin directory: %v\n", err)
		return ""
	}

	//Mount all subfolders in mountPathOrigin that start with @ to mountPathDest without the @, only those in the first level
	// Get all entries in the mountPathOrigin directory
	entries, err := os.ReadDir(mountPathOrigin)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return ""
	}

	// Iterate through entries and mount those starting with @
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "@") {
			// Create destination directory without @ prefix
			destPath := filepath.Join(mountPathDest, strings.TrimPrefix(entry.Name(), "@"))
			if err := os.MkdirAll(destPath, 0744); err != nil {
				fmt.Printf("Error creating destination directory: %v\n", err)
				os.Exit(1)
				continue
			}

			// Mount the subvolume
			sourcePath := filepath.Join(mountPathOrigin, entry.Name())
			if err := cmdx.RunCommandPrintOutput("mount", "--bind", sourcePath, destPath); err != nil {
				fmt.Printf("Error mounting directory %s to %s: %v\n", sourcePath, destPath, err)
				os.Exit(1)
				continue
			}
		}
	}

	return mountPathDest

}

func unmountFalseMountpoint(mountPathOrigin string, mountPathDest string) error {
	// Get all entries in the destination directory
	entries, err := os.ReadDir(mountPathDest)
	if err != nil {
		return fmt.Errorf("error reading directory: %v", err)
	}

	// Unmount all directories in the destination path
	for _, entry := range entries {
		if entry.IsDir() {
			destPath := filepath.Join(mountPathDest, entry.Name())

			// Unmount the bind mount
			if err := cmdx.RunCommandPrintOutput("umount", destPath); err != nil {
				return fmt.Errorf("error unmounting directory %s: %v", destPath, err)
			}

			// Remove the empty directory
			if err := os.Remove(destPath); err != nil {
				return fmt.Errorf("error removing directory %s: %v", destPath, err)
			}
		}
	}

	// Remove the parent destination directory
	if err := os.Remove(mountPathDest); err != nil {
		return fmt.Errorf("error removing destination directory %s: %v", mountPathDest, err)
	}

	return nil
}

var MountPath string
var SubvolPath1 string

func TestMain(m *testing.M) {

	tempfilename := createBtrfsImage()
	// Create a temporary mount directory
	testMountPath := filepath.Join(os.TempDir(), "testbtrfs")

	createTestBtrfsMount(tempfilename)

	testSubvolPath := createTestSubvolume(testMountPath, "@testsubvol")

	MountPath = createFalseMountpoint(testMountPath, filepath.Join(os.TempDir(), "falseroot"))
	SubvolPath1 = filepath.Join(MountPath, "testsubvol")
	// Run the tests
	code := m.Run()

	// Unmount the Btrfs filesystem
	if err := cmdx.RunCommandPrintOutput("umount", testMountPath); err != nil {
		fmt.Printf("Error unmounting Btrfs filesystem: %v\n", err)
	}

	unmountFalseMountpoint(testMountPath, filepath.Join(os.TempDir(), "falseroot"))
	// Clean up
	btrfs.SubvolDelete(testSubvolPath)
	os.RemoveAll(testMountPath)
	os.Remove(tempfilename)

	os.Exit(code)
}
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

	// Clean up test subvolume
	if err := os.RemoveAll(testSubvolPath); err != nil {
		t.Logf("Warning: Failed to cleanup test subvolume: %v", err)
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
