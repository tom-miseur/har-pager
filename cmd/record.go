/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"time"
	"encoding/json"

	"github.com/spf13/cobra"
	"github.com/manifoldco/promptui"
)

var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "Record a new series of pages",
	Long: `Record a new series of one or more pages.
Make sure this runs in the background while recording a HAR file elsewhere.
	
Usage:
har-pager record MyUserJourney
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0] + ".log"

		if _, err := os.Stat(filename); os.IsNotExist(err) {
			fmt.Println("Recording started. To finish the recording, press CTRL/CMD+C")
		} else {
			fmt.Println(filename + " already exists. Overwrite?")

			prompt := promptui.Select{
        Label: "Select[Yes/No]",
        Items: []string{"Yes", "No"},
    }
    _, result, err := prompt.Run()
    if err != nil {
				fmt.Println("Prompt failed %v\n", err)
    }
			if result == "No" {
				os.Exit(0)
			} else {
				fmt.Println(filename + " will be overwritten. To finish the recording, press CTRL/CMD+C")
			}
		}

		// main page name prompt loop
		pages := promptLoop()

		// done (via user-initiated CTRL/CMD+C)
		fmt.Println("Save recording?")

		prompt := promptui.Select{
        Label: "Select[Yes/No]",
        Items: []string{"Yes", "No"},
    }
    _, result, err := prompt.Run()
    if err != nil {
				fmt.Println("Prompt failed %v\n", err)
    }
	
		if result == "No" {
			os.Exit(0)
		} else {
			err := savePagesToJSON(pages, filename)

			if err != nil {
				fmt.Println("Error saving pages to JSON file: ", err)
				os.Exit(1)
			}
		}

    fmt.Println("Done! Recording saved to " + filename)	
		
    os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(recordCmd)
}

func promptLoop() []PageLog {
	pages := []PageLog{}

	for {
		// TODO: validate the input
		// validate := func(input string) error {
		// 	_, err := strconv.ParseFloat(input, 64)
		// 	if err != nil {
		// 		return errors.New("Invalid number")
		// 	}
		// 	return nil
		// }
	
		prompt := promptui.Prompt{
			Label:    "Page name",
			// ValiStarted: validate,
			HideEntered: true,
		}
	
		result, err := prompt.Run()
	
		if err != nil {
			// CTRL/CMD+C
			fmt.Printf("Recording complete. %v\n", err)
			
			fmt.Println("Pages captured: ", len(pages))
			return pages
		}
	
		// "end" the previous page (just a visual cue for the user; doesn't change the JSON output)
		if len(pages) > 0 {
			fmt.Printf("Ending page %s. Page duration: %s\n", pages[len(pages)-1].Name, time.Since(pages[len(pages)-1].Started))
		}

		fmt.Printf("Starting page '%s'. Enter next page, or CTRL/CMD+C to stop recording.\n", result)

		// capture the time
		n := time.Now()

		page := PageLog{
			Name: result,
			Started: n,
		}
		
		pages = append(pages, page)
	}
}

func savePagesToJSON(pages []PageLog, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(pages)
	if err != nil {
		return err
	}

	return nil
}