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
		interface_create_subvolume(cmd, args, btrfs.SubvolCreate, CreateSubvolumeRedhatStyle)
	},
}

// func(dir interface{}) (interface{}, error) { return check_btrfs_subvolume(dir.(string)) }
func interface_create_subvolume(cmd *cobra.Command, args []string, create_subvol func(string) error, create_redhat_subvol func(string) error) {
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
	// TODO mount new subvolume into destiny
}

func CreateSubvolumeRedhatStyle(subvolume string) error {
	// Get Physical device to perform operations
	var subvolume_derived string

	if _, err := os.Stat(subvolume); os.IsExist(err) {
		return fmt.Errorf("subvolume already exists: %s", subvolume)
	}
	subvolume_derived = filepath.Dir(subvolume)
	subvolume_basename := filepath.Base(subvolume)
	phDev, _, err := GetMountpoint(subvolume_derived)
	if err != nil {
		return fmt.Errorf("failed to get mountpoint: %w", err)
	}

	// Create and mount temporary directory
	tempMountPath, err := createAndMountTempDir(phDev)
	if err != nil {
		return err
	}
	defer cleanupTempMount(tempMountPath)

	// Create subvolume
	err = btrfs.SubvolCreate(tempMountPath + "/@" + subvolume_basename)
	//err = createSubvolumeAtPath(subvolume, tempMountPath)
	if err != nil {
		return err
	}

	return nil
}

func createAndMountTempDir(phDev string) (string, error) {
	mountUUID := uuid.New().String()
	tempMountPath := filepath.Join(tmppath, "/mount-create-", mountUUID)

	err := os.MkdirAll(tempMountPath, 0744)
	if err != nil {
		return "", fmt.Errorf("failed to create temporary mount directory: %w", err)
	}

	args := []string{phDev, tempMountPath}
	_, err = cmdx.RunCommandReturnOutputWithDirAndEnv("mount", tmppath, nil, args...)
	if err != nil {
		os.RemoveAll(tempMountPath)
		return "", fmt.Errorf("failed to mount device: %w", err)
	}

	return tempMountPath, nil
}

func cleanupTempMount(tempMountPath string) error {
	args := []string{tempMountPath}
	_, err := cmdx.RunCommandReturnOutputWithDirAndEnv("umount", tmppath, nil, args...)
	if err != nil {
		return fmt.Errorf("failed to unmount device: %w", err)
	}

	err = os.RemoveAll(tempMountPath)
	if err != nil {
		return fmt.Errorf("failed to remove temporary mount directory: %w", err)
	}
	return nil
}

func createSubvolumeAtPath(localsubvolume string, tempMountPath string) error {
	basename := filepath.Base(localsubvolume)
	newSubvolPath := filepath.Join(tempMountPath, basename)

	err := btrfs.SubvolCreate(newSubvolPath)
	if err != nil {
		return fmt.Errorf("failed to create subvolume: %w", err)
	}
	return nil
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
