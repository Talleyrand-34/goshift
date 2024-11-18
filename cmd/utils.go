package cmd

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/containerd/btrfs/v2"
)

// Check_btrfs_subvolume This function checks if a string is a valid btrfs subvolume in the host
func check_btrfs_subvolume(dirstring string) (finaldir string, err error) {
	dir := ""
	if filepath.IsAbs(dirstring) {
		dir = dirstring
		// return dirstring
	} else {
		dir, _ = os.Getwd()
		dir += dirstring
		// return dir + dirstring
	}
	// check if it is subvolume
	if btrfs.IsSubvolume(dir) != nil {
		return dir, errors.New(dir + "is not a subvolume")
	}
	return dir, nil
}

// Prints formatted error message to the user
// func errorMessage(userinput string, err error) {
// 	log.Fatal("Error: ", err.Error(), "|| Input was ", userinput)
// }

// func findCommand(command string, subcommand string) (*Command, error) {
// 	for _, cmd := range Commands {
// 		if cmd.CommL1 == command && cmd.CommL2 == subcommand {
// 			return &cmd, nil
// 		}
// 	}
// 	return nil, fmt.Errorf("command not found: %s %s", command, subcommand)
// }
//
// func printUsage() {
// 	fmt.Println("Usage: program <command> <subcommand> <route> [additional args...]")
// 	fmt.Println("Available commands:")
// 	for _, cmd := range Commands {
// 		fmt.Printf("  %s %s: %s\n", cmd.CommL1, cmd.CommL2, cmd.Description)
// 	}
// }
