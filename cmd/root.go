package cmd

import (
	"os"
	"time"

	"github.com/spf13/cobra"
)

type Har struct {
	Log Log `json:"log"`
}

type Log struct {
	Version string `json:"version"`
	Creator map[string]interface{} `json:"creator"`
	Pages []Page `json:"pages"`
	Entries []Entry `json:"entries"`
	Comment string `json:"comment"`
}

type Page struct {
	StartedDateTime time.Time `json:"startedDateTime"`
	ID string `json:"id"`
	Title string `json:"title"`
	PageTimings PageTimings `json:"pageTimings"`
	Comment string `json:"comment"`
}

type Entry struct {
	Pageref string `json:"pageref"`
	StartedDateTime time.Time `json:"startedDateTime"`
	Time int `json:"time"`
	Request map[string]interface{} `json:"request"`
	Response map[string]interface{} `json:"response"`
	Cache map[string]interface{} `json:"cache"`
	Timings map[string]interface{} `json:"timings"`
	ServerIPAddress string `json:"serverIPAddress"`
	Connection string `json:"connection"`
	Comment string `json:"comment"`	
}

type PageTimings struct {
	OnContentLoad int `json:"onContentLoad"`
	OnLoad int `json:"onLoad"`
	Comment string `json:"comment"`
}

type PageLog struct {
	Name string `json:"name"`
	Started time.Time `json:"started"`
}

var rootCmd = &cobra.Command{
	Use:   "har-pager",
	Short: "har-pager is a tool to record and merge pages into a HAR file",
	Long: `har-pager is a tool to record and merge pages into a HAR file,
where the HAR file is generated elsewhere (such as with Chrome's DevTools,
or a web proxy like Fiddler or Charlesproxy).

har-pager has two modes of operation: record and merge.

In record mode, you enter the names of the pages you want to generate.
Typically, these pages represent higher-level "user actions", such as
"Login", or "Add to Cart". After entering a page name, perform the action
that results in HTTP requests being captured into the HAR file. Continue
this until you have created all the pages you want, at which point you can
stop recording with CTRL/CMD+C to save the page log.

In merge mode, the recorded page log is merged with the HAR file, using
timestamps of when the pages were created to determine which requests 
belong to which page.

Usage:

har-pager record MyPages
har-pager merge MyPages MyRecording.har MyRecording-merged.har

For more information, see https://github.com/tom-miseur/har-pager
`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
