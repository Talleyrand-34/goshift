/*
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package main

import (
	// "flag"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	// "text/tabwriter"
	"github.com/containerd/btrfs/v2"
	// tea "github.com/charmbracelet/bubbletea"
)

func manageSubvolume(args []string) {
	// check run conditions
	if len(args) < 1 {
		log.Fatal("Subvolume needs a command", os.Args[1])
		os.Exit(1)
		return
	}
	if len(args) < 2 {
		log.Fatal("create needs a valid folder", os.Args[1])
		os.Exit(1)
	}
	// Path treatment
	dir := ""
	if filepath.IsAbs(args[1]) {
		dir = args[1]
	} else {
		dir, _ = os.Getwd()
		dir = dir + args[1]
	}
	fmt.Println(dir)

	fmt.Println("Is a valid subvolume:", btrfs.IsSubvolume(dir))
	// command section
	switch args[0] {
	case "list":
		treelist, err := btrfs.SubvolList(dir)
		if err != nil {
			log.Fatalln(err)
		}
		jsonData, err := transformInfoArray(treelist)
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return
		}

		fmt.Println(string(jsonData))
		// fmt.Printf("%#v\n", treelist)
		// for _, sv := range treelist {
		// 	fmt.Println(sv)
		// }

	case "show":
		info, err := btrfs.SubvolInfo(dir)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%#v\n", info)

	case "create":
		// if (){

		err := btrfs.SubvolCreate(dir)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("Subvolume needs a command", os.Args[1])
		os.Exit(1)

		// }
	}
}

func transformInfoArray(infos []btrfs.Info) ([]byte, error) {
	return json.Marshal(infos)
}

func main() {
	args := os.Args[1:]
	if args[0] == "" {
		log.Fatal("unknown command", os.Args[1])
		os.Exit(1)
	}
	switch args[0] {
	case "subvolume":
		fmt.Println("case subvolume")
		manageSubvolume(args[1:])

	case "snapshot":
		fmt.Println("case snapshot")
	default:
		log.Fatal("unknown command", os.Args[1])
		os.Exit(1)
	}
	os.Exit(0)
}

// p := tea.NewProgram(initialModel())
// if _, err := p.Run(); err != nil {
//     fmt.Printf("Alas, there's been an error: %v", err)
//     os.Exit(1)
// }

// switch os.Args[1] {
// case "create":
// 	if err := btrfs.SubvolCreate(os.Args[2]); err != nil {
// 		log.Fatalln(err)
// 	}
// case "snapshot":
// 	if err := btrfs.SubvolSnapshot(os.Args[3], os.Args[2], readonly); err != nil {
// 		log.Fatalln(err)
// 	}
// case "delete":
// 	if err := btrfs.SubvolDelete(os.Args[2]); err != nil {
// 		log.Fatalln(err)
// 	}
// case "list":
// 	infos, err := btrfs.SubvolList(os.Args[2])
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	tw := tabwriter.NewWriter(os.Stdout, 0, 8, 4, '\t', 0)
//
// 	fmt.Fprintf(tw, "ID\tParent\tTopLevel\tGen\tOGen\tUUID\tParentUUID\tPath\n")
//
// 	for _, subvol := range infos {
// 		fmt.Fprintf(tw, "%d\t%d\t%d\t%d\t%d\t%s\t%s\t%s\n",
// 			subvol.ID, subvol.ParentID, subvol.TopLevelID,
// 			subvol.Generation, subvol.OriginalGeneration, subvol.UUID, subvol.ParentUUID,
// 			subvol.Path)
//
// 	}
//
// 	tw.Flush()
// case "show":
// 	info, err := btrfs.SubvolInfo(os.Args[2])
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
//
// 	fmt.Printf("%#v\n", info)
// default:
// 	log.Fatal("unknown command", os.Args[1])
// }

// TUI
//
// type model struct {
//     choices  []string           // items on the to-do list
//     cursor   int                // which to-do list item our cursor is pointing at
//     selected map[int]struct{}   // which to-do items are selected
// }
// func initialModel() model {
// 	return model{
// 		// Our to-do list is a grocery list
// 		choices:  []string{"List Subvolumes", "Buy celery", "Buy kohlrabi"},
//
// 		// A map which indicates which choices are selected. We're using
// 		// the  map like a mathematical set. The keys refer to the indexes
// 		// of the `choices` slice, above.
// 		selected: make(map[int]struct{}),
// 	}
// }
// func (m model) Init() tea.Cmd {
//     // Just return `nil`, which means "no I/O right now, please."
//     return nil
// }
//
// func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
//     switch msg := msg.(type) {
//
//     // Is it a key press?
//     case tea.KeyMsg:
//
//         // Cool, what was the actual key pressed?
//         switch msg.String() {
//
//         // These keys should exit the program.
//         case "ctrl+c", "q":
//             return m, tea.Quit
//
//         // The "up" and "k" keys move the cursor up
//         case "up", "k":
//             if m.cursor > 0 {
//                 m.cursor--
//             }
//
//         // The "down" and "j" keys move the cursor down
//         case "down", "j":
//             if m.cursor < len(m.choices)-1 {
//                 m.cursor++
//             }
//
//         // The "enter" key and the spacebar (a literal space) toggle
//         // the selected state for the item that the cursor is pointing at.
//         case "enter", " ":
//             _, ok := m.selected[m.cursor]
//             if ok {
//                 delete(m.selected, m.cursor)
//             } else {
//                 m.selected[m.cursor] = struct{}{}
//             }
//         }
//     }
//
//     // Return the updated model to the Bubble Tea runtime for processing.
//     // Note that we're not returning a command.
//     return m, nil
// }
// func (m model) View() string {
//     // The header
//     s := "What should we buy at the market?\n\n"
//
//     // Iterate over our choices
//     for i, choice := range m.choices {
//
//         // Is the cursor pointing at this choice?
//         cursor := " " // no cursor
//         if m.cursor == i {
//             cursor = ">" // cursor!
//         }
//
//         // Is this choice selected?
//         checked := " " // not selected
//         if _, ok := m.selected[i]; ok {
//             checked = "x" // selected!
//         }
//
//         // Render the row
//         s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
//     }
//
//     // The footer
//     s += "\nPress q to quit.\n"
//
//     // Send the UI for rendering
//     return s
// }
