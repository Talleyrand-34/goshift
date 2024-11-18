package cmd_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/containerd/btrfs/v2"
	"github.com/kgs19/cmdx"
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
func CreateTestSubfolder(mountPath string, name string) string {
	testSubvolPath := filepath.Join(mountPath, name)
	if err := os.MkdirAll(testSubvolPath, 0744); err != nil {
		fmt.Printf("Error creating test subfolder: %v\n", err)
		os.Exit(1)
	}
	return testSubvolPath
}

var MountPath string
var SubvolPath1 string
var RootPath string
var SubFolderPath string

func TestMain(m *testing.M) {

	tempfilename := createBtrfsImage()
	// Create a temporary mount directory
	testMountPath := filepath.Join(os.TempDir(), "testbtrfs")
	RootPath = testMountPath

	createTestBtrfsMount(tempfilename)

	testSubvolPath := createTestSubvolume(testMountPath, "@testsubvol")

	MountPath = createFalseMountpoint(testMountPath, filepath.Join(os.TempDir(), "falseroot"))
	SubvolPath1 = filepath.Join(MountPath, "testsubvol")
	SubFolderPath = CreateTestSubfolder(SubvolPath1, "testsubfolder")
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
