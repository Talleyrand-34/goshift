/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/containerd/btrfs/v2"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		route := args[len(args)-1]
		if btrfs_sv, err := check_btrfs_subvolume(route); err != nil {
			btrfs_subfolder, err := isPathOnBtrfs(btrfs_sv)
			print(btrfs_subfolder, "\n")
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			if btrfs_subfolder {
				print(btrfs_sv, " is a subfolder\n")
			}
			log.Fatal(btrfs_sv, " is not a valid btrfs subvolume\n")
			return
		}
		var info btrfs.Info
		info, err := btrfs.SubvolInfo(route)
		if err != nil {
			log.Fatal("Error obtaining the info of ", route, "; Not enough priviledges")
			return
		}
		jsoninfo, err := json.Marshal(info)
		if err != nil {
			log.Fatal("error marshaling JSON:", err)
		}
		fmt.Print(string(jsoninfo))
	},
}

func isPathOnBtrfs(path string) (bool, error) {
	var statfs unix.Statfs_t
	if err := unix.Statfs(path, &statfs); err != nil {
		return false, err
	}
	return statfs.Type == unix.BTRFS_SUPER_MAGIC, nil
}

func init() {
	subvolumeCmd.AddCommand(infoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
