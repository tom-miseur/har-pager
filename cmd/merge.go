/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/spf13/cobra"
)

var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge a page log into a HAR file",
	Long: `The merge command is the second step in the har-pager process.
Combine a page log, created using the record command, with a HAR file
that was recorded at the same time as the page log.

The last parameter represents the output file name containing the merged data.

Usage:
har-pager merge MyUserJourney MyUserJourney.har MyUserJourney-merged.har
`,
	Run: func(cmd *cobra.Command, args []string) {
		pages := readPagesFromJSON(args[0] + ".log")
		fmt.Println(pages)

		har := readHAR(args[1])

		merged := mergePagesIntoHAR(pages, har)

		saveHARToFile(merged, args[2])
	},
}

func init() {
	rootCmd.AddCommand(mergeCmd)
}

func readPagesFromJSON(filename string) []PageLog {
	pagesJSON, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Println("Error reading JSON file: ", err)
		os.Exit(1)
	}

	var pages []PageLog
	err = json.Unmarshal(pagesJSON, &pages)
	if err != nil {
		fmt.Println("Error unmarshalling JSON: ", err)
	}

	return pages
}

func readHAR(filename string) Har {
	harJSON, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Println("Error reading HAR file: ", err)
		os.Exit(1)
	}

	// handle possible BOM
	harJSON = bytes.TrimPrefix(harJSON, []byte("\xef\xbb\xbf"))
	
	var har Har
	err = json.Unmarshal(harJSON, &har)

	if err != nil {
		fmt.Println("Error unmarshalling JSON: ", err)
		os.Exit(1)
	}

	return har
}

func mergePagesIntoHAR(pageLogs []PageLog, har Har) Har {
	var harPages []Page
	
	for _, page := range pageLogs {
		name := page.Name
		started := page.Started

		page := Page{
			StartedDateTime: started,
			ID: name,
			Title: "",
			PageTimings: PageTimings{
				OnContentLoad: -1,
				OnLoad: -1,
				Comment: "",
			},
			Comment: "Added by har-pager",
		}

		harPages = append(harPages, page)
	}

	har.Log.Pages = harPages

	for i, entry := range har.Log.Entries {
		started := entry.StartedDateTime

		for j, page := range pageLogs {
			pageStarted := page.Started

			if started.After(pageStarted) && (j == len(pageLogs)-1 || started.Before(pageLogs[j+1].Started)) {
				har.Log.Entries[i].Pageref = page.Name
			}
		}
	}

	return har
}

func saveHARToFile(har Har, filename string) {
	harJSON, err := json.Marshal(har)
	if err != nil {
		fmt.Println("Error marshalling HAR to JSON: ", err)
	}

	err = ioutil.WriteFile(filename, harJSON, 0644)
	if err != nil {
		fmt.Println("Error writing HAR to file: ", err)
	}
}