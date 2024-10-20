// Package This package provides a btrfs volumes management tool
package main

import (
	// "flag"

	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/containerd/btrfs/v2"
	// tea "github.com/charmbracelet/bubbletea"
	// "text/tabwriter"
)

type Command struct {
	CommL1        string
	CommL2        string
	Description   string
	Execute       func(interface{}) (interface{}, error)
	Check_path    func(interface{}) (interface{}, error)
	Error_msg     func(string, error)
	Validate_args func(interface{}) (interface{}, error)
	Json_msg      func(interface{}) (interface{}, error)
}

var Commands = []Command{
	{
		CommL1:        "subvolume",
		CommL2:        "create",
		Description:   "Creates a subvolume in btrfs volume root",
		Execute:       func(dir interface{}) (interface{}, error) { return btrfs.SubvolList(dir.(string)) },
		Check_path:    func(dir interface{}) (interface{}, error) { return check_btrfs_subvolume(dir.(string)) },
		Error_msg:     errorMessage,
		Validate_args: nil,
		Json_msg:      func(btrfs_struct interface{}) (interface{}, error) { return json.Marshal(btrfs_struct.(btrfs.Info)) },
	},
	{
		CommL1:        "subvolume",
		CommL2:        "info",
		Description:   "Shows subvolume info",
		Execute:       func(dir interface{}) (interface{}, error) { return btrfs.SubvolInfo(dir.(string)) },
		Check_path:    func(dir interface{}) (interface{}, error) { return check_btrfs_subvolume(dir.(string)) },
		Error_msg:     errorMessage,
		Validate_args: nil,
		Json_msg:      func(btrfs_struct interface{}) (interface{}, error) { return json.Marshal(btrfs_struct.(btrfs.Info)) },
	},
}

func BtrfsCommandExecuter(
	command string,
	subcommand string,
	route string,
	args []string,
) (jsondata string, err error) {
	cmd, err := findCommand(command, subcommand)
	if err != nil {
		return "", fmt.Errorf("command not found: %v", err)
	}
	// // Check if the route is valid
	// checkResult, checkErr := cmd.Check_path(route)
	// if checkErr != nil {
	// 	return "", fmt.Errorf("Error checking route: %v", checkErr)
	// }
	// if !checkResult.(bool) {
	// 	return "", fmt.Errorf("Invalid route: %v", route)
	// }

	btrfsfn, err := cmd.Execute(route)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// Process the btrfsfn through Json_msg function
	jsonResult, err := cmd.Json_msg(btrfsfn)
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %v", err)
	}

	// Convert jsonResult to string
	jsondata = string(jsonResult.([]byte))

	return jsondata, nil
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		printUsage()
		os.Exit(0)
	}

	if len(args) < 3 {
		log.Fatal("Insufficient arguments. Use 'help' for usage information.")
	}

	command := args[0]
	subcommand := args[1]
	route := args[2]
	var additionalArgs []string

	if len(args) > 3 {
		additionalArgs = args[3:]
	} else {
		additionalArgs = []string{""}
	}

	fmt.Printf("Executing: %s %s on route %s with args: %v\n", command, subcommand, route, additionalArgs)

	info, err := BtrfsCommandExecuter(command, subcommand, route, additionalArgs)
	if err != nil {
		log.Fatalf("Error executing command: %v", err)
	}

	fmt.Println("Command output:")
	fmt.Println(info)

	os.Exit(0)
}

func printUsage() {
	fmt.Println("Usage: program <command> <subcommand> <route> [additional args...]")
	fmt.Println("Available commands:")
	for _, cmd := range Commands {
		fmt.Printf("  %s %s: %s\n", cmd.CommL1, cmd.CommL2, cmd.Description)
	}
}
