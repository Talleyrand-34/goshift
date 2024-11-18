// Package This package provides a btrfs volumes management tool
package main

// "flag"

// tea "github.com/charmbracelet/bubbletea"
// "text/tabwriter"

// var config_file str

// import (
// 	"goshift/cmd"
// )
//
// func main() {
// 	cmd.Execute()
// }

// func main() {
// 	args := os.Args[1:]
// 	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
// 		printUsage()
// 		os.Exit(0)
// 	}
//
// 	if len(args) < 3 {
// 		log.Fatal("Insufficient arguments. Use 'help' for usage information.")
// 	}
//
// 	command := args[0]
// 	subcommand := args[1]
// 	route := args[2]
// 	var additionalArgs []string
//
// 	if len(args) > 3 {
// 		additionalArgs = args[3:]
// 	} else {
// 		additionalArgs = []string{""}
// 	}
//
// 	fmt.Printf("Executing: %s %s on route %s with args: %v\n", command, subcommand, route, additionalArgs)
//
// 	info, err := BtrfsCommandExecuter(command, subcommand, route, additionalArgs)
// 	if err != nil {
// 		log.Fatalf("Error executing command: %v", err)
// 	}
//
// 	fmt.Println("Command output:")
// 	fmt.Println(info)
//
// 	os.Exit(0)
// }
