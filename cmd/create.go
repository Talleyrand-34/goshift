/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/containerd/btrfs/v2"
	"github.com/google/uuid"
	"github.com/kgs19/cmdx"
	"github.com/spf13/cobra"
)

var (
	redhatstd bool
	tmppath   string = "/tmp/goshift"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create a subvolume",
	Long: `

	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		interface_create_subvolume(cmd, args, btrfs.SubvolCreate, create_subvolume_redhat_style)
	},
}

// func(dir interface{}) (interface{}, error) { return check_btrfs_subvolume(dir.(string)) }
func interface_create_subvolume(cmd *cobra.Command, args []string, create_subvol func(string) error, create_redhat_subvol func(string)) {
	redhatstd, _ := cmd.Flags().GetBool("redhatstd")
	route := args[len(args)-1]
	if !redhatstd {
		// err := btrfs.SubvolCreate(route)
		err := create_subvol(route)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	create_redhat_subvol(route)
}

func create_subvolume_redhat_style(subvolume string) {
	// func create_subvolume_redhat_style(subvolume string) {
	// Get the base name (file name or last directory in the path)
	// baseName := filepath.Base(subvolume)
	// Get the directory name (path without the last element)
	// dirName := filepath.Dir(subvolume)
	// TODO err management
	// Get Physical device to perform operations
	phDev, psubv, _ := GetMountpoint(subvolume)
	print(phDev, "\n")
	print(psubv, "\n")
	// Mount device on "/tmp/goshift/mount-create + uuid"
	// Generate a unique UUID for the temporary mount point
	mountUUID := uuid.New().String()
	tempMountPath := filepath.Join(tmppath, "/mount-create-", mountUUID)
	// Create the temporary mount directory
	// err := cmdx.RunCommandPrintOutput("mkdir", "-p", tempMountPath)
	err := os.MkdirAll(tempMountPath, 0744)
	if err != nil {
		fmt.Printf("Error creating temporary mount directory: %v\n", err)
		return
	}
	// Mount device on temporary path
	args := []string{phDev, tempMountPath}
	_, err = cmdx.RunCommandReturnOutputWithDirAndEnv("mount", tmppath, nil, args...)
	if err != nil {
		fmt.Printf("Error mounting device: %v\n", err)
		return
	}
	// Create subvolume
	basename_subvolume := filepath.Base(subvolume)
	newSubvolPath := filepath.Join(tempMountPath, basename_subvolume)
	err = btrfs.SubvolCreate(newSubvolPath)
	if err != nil {
		fmt.Printf("Error creating subvolume")
	}

	// Unmount device
	// unmountCmd := exec.Command("umount", tempMountPath)
	// err = unmountCmd.Run()
	args = []string{tempMountPath}
	_, err = cmdx.RunCommandReturnOutputWithDirAndEnv("umount", tmppath, nil, args...)
	if err != nil {
		fmt.Printf("Error unmounting device: %v\n", err)
	}

	// Clean up: remove the temporary mount directory
	err = os.RemoveAll(tempMountPath)
	if err != nil {
		fmt.Printf("Error removing temporary mount directory: %v\n", err)
	}
}

func GetMountpoint(subvolume string) (string, string, error) {
	//args := []string{"--target=", subvolume}
	args := []string{"--target=" + subvolume}
	cmd, err := cmdx.RunCommandReturnOutputWithDirAndEnv("findmnt", "/tmp", nil, args...)
	if err != nil {
		return "", "", err
	} // Trim whitespace and newlines from the output
	arrayfindmnt := strings.Fields(string(cmd))
	phDev, subv := SplitStringWithBrackets(arrayfindmnt[5])
	return phDev, subv, nil
	// return cmd, cmd, nil
}

func SplitStringWithBrackets(input string) (string, string) {
	// Find the index of the opening bracket
	openBracketIndex := strings.Index(input, "[")

	if openBracketIndex == -1 {
		// If there's no opening bracket, return the whole string and an empty substring
		return input, ""
	}

	// Split the string at the opening bracket
	mainString := input[:openBracketIndex]
	substring := input[openBracketIndex:]

	return mainString, substring
}

func init() {
	subvolumeCmd.AddCommand(createCmd)
	subvolumeCmd.PersistentFlags().BoolVarP(&redhatstd, "redhatstd", "s", false, "Use redhat creation style")
}
