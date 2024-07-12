# har-pager

har-pager is a CLI utility that allows HAR [pages](http://www.softwareishard.com/blog/har-12-spec/#pages) to be recorded in such a way that they can subsequently be merged into a HAR file, captured elsewhere at the same time as the recording, that does not contain them.

### What is a HAR page?

The HAR spec caters for associating `requests` to a parent `page` object. Although optional, it can be useful as a means of grouping requests that share something in common. The original intent behind a HAR page was likely for it to represent a browser navigation, it can also be used to denote the "user action" that the requests belong to. For example, a "Submit Login" page could contain all the requests that took place when the user submitted a login form.

### Use case with k6

[k6](https://k6.io) uses HAR page information to generate [groups](https://grafana.com/docs/k6/latest/using-k6/tags-and-groups/#groups) when converting HAR recordings to k6 scripts using the Import HAR functionality in Grafana Cloud k6. The main benefit of this is to make the resulting auto-generated HTTP code easier to interpret, particularly when dealing with lengthy flows with many requests.

### Installation

Make sure you have Go installed on your machine. If not, you can download it from [here](https://golang.org/dl/).

From source:
1. Clone the repository
2. Run `go build` in the root directory
3. Move the resulting `har-pager` binary to a location in your PATH environment variable (so that it can be run from anywhere)

Pre-built binaries:

Coming soon.

### Usage

har-pager has two modes of operation: `record` and `merge`.

#### Record mode

```shell
har-pager record MyPages
```

In record mode, you enter the names of the pages you want to generate. Running the above command will start the recording process:

```shell
Recording started. To finish the recording, press CTRL/CMD+C
v Page name: 
```

Now would be a good time to also start the HAR recording. This could be done through Chrome DevTools' Network tab, a web proxy like Telerik's [Fiddler](https://www.telerik.com/download/fiddler) or [mitmproxy](https://mitmproxy.org/), or even the [k6 Browser Recorder](https://grafana.com/docs/k6/latest/using-k6/test-authoring/create-tests-from-recordings/using-the-browser-recorder/).

Enter a page name, hit return, and then perform the action that results in HTTP requests. For example:

```shell
v Page name: Navigate to Homepage
Starting page 'Navigate to Homepage'. Enter next page, or CTRL/CMD+C to stop recording.
```

Continue this loop as necessary, using CTRL/CMD+C to stop the recording and save the page log.

#### Merge mode

```shell
har-pager merge MyPages MyRecording.har MyRecording-merged.har
```

In merge mode, a recorded page log can be merged with the corresponding HAR file (2nd parameter). The 3rd parameter is the output HAR file.

The merging process is quite simple: the generated page log contains a timestamp for each page that corresponds to when it was entered during the recording. This timestamp is then compared with `request` `startDateTime` timestamps in the HAR file to determine whether a request took place between this timestamp and that of the subsequent page. Should that be the case, a `pageRef` field is added to the request object, linking it to the page.
