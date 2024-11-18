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

func createTestBtrfsMount(imagePath string, localmountPath string) {
	// Create a temporary mount directory
	if err := os.MkdirAll(localmountPath, 0744); err != nil {
		fmt.Printf("Error creating temporary mount directory: %v\n", err)
		os.Exit(1)
	}

	// Mount the temporary Btrfs filesystem
	if err := cmdx.RunCommandPrintOutput("mount", imagePath, localmountPath); err != nil {
		fmt.Printf("Error mounting Btrfs filesystem on %s: %v\n", localmountPath, err)
		os.Exit(1)
	}
}

func createTestSubvolume(subvolPath string) {

	// Create test subvolume
	if err := btrfs.SubvolCreate(subvolPath); err != nil {
		fmt.Printf("Failed to create test subvolume %s: %v", subvolPath, err)
		os.Exit(1)
	}

}

func createFalseMountpoint(mountPathOrigin string, mountPathDest string) {
	if err := os.MkdirAll(mountPathDest, 0744); err != nil {
		fmt.Printf("Error creating mountPathOrigin directory: %v\n", err)
	}
	rootsourcePath := mountPathOrigin + "/@"
	if err := btrfs.SubvolCreate(mountPathOrigin + "/@"); err != nil {
		fmt.Printf("Error creating root subvolume %s: %v", mountPathOrigin+"/@", err)
	}
	if err := cmdx.RunCommandPrintOutput("mount", "--bind", rootsourcePath, mountPathDest); err != nil {
		fmt.Printf("Error mounting directory %s to %s: %v\n", rootsourcePath, mountPathDest, err)
		os.Exit(1)
	}

	//Mount all subfolders in mountPathOrigin that start with @ to mountPathDest without the @, only those in the first level
	// Get all entries in the mountPathOrigin directory
	entries, err := os.ReadDir(mountPathOrigin)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
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
func CreateTestSubfolder(Path string) {
	if err := os.MkdirAll(Path, 0744); err != nil {
		fmt.Printf("Error creating test subfolder: %v\n", err)
		os.Exit(1)
	}
}

var RootPath string       // /tmp/testbtrfs
var RootSubvolPath string // /tmp/testbtrfs/@testsubvol
var MountPath string      // /tmp/falseroot
var SubvolPath1 string    // /tmp/falseroot/testsubvol
var SubFolderPath string  // /tmp/falseroot/testsubfolder

func TestMain(m *testing.M) {

	tempfilename := createBtrfsImage()
	// Create a temporary mount directory
	RootPath = filepath.Join(os.TempDir(), "testbtrfs")
	createTestBtrfsMount(tempfilename, RootPath)
	// create subvolume in root
	RootSubvolPath = RootPath + "/@testsubvol"
	createTestSubvolume(RootSubvolPath)

	// create false mountpoint
	MountPath = filepath.Join(os.TempDir(), "falseroot")
	createFalseMountpoint(RootPath, MountPath)
	// create mountpoint for subvolume in false mountpoint
	SubvolPath1 = filepath.Join(MountPath, "testsubvol")
	SubFolderPath = filepath.Join(SubvolPath1, "testsubfolder")
	CreateTestSubfolder(SubFolderPath)
	// Run the tests
	code := m.Run()

	// Unmount the Btrfs filesystem
	if err := cmdx.RunCommandPrintOutput("umount", RootPath); err != nil {
		fmt.Printf("Error unmounting Btrfs filesystem: %v\n", err)
	}

	// Clean up
	unmountFalseMountpoint(RootPath, filepath.Join(os.TempDir(), "falseroot"))
	btrfs.SubvolDelete(RootSubvolPath)
	os.RemoveAll(RootPath)
	os.Remove(tempfilename)

	os.Exit(code)
}
