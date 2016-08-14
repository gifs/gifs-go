# gifs go

[![GoDoc](https://godoc.org/github.com/gifs/gifs-go?status.svg)](https://godoc.org/github.com/gifs/gifs-go)

Go package for interacting with the gifs.com API for transcoding and importing media to `.mp4`, `.jpg`, `.gif` outputs. Import all your videos by uploading or passing in any source. It is super easy to import media and integrate in your code.

## Examples

Please see the examples below on how to import a single source, then in bulk each in one call.

You can also see more usage and convenience methods in files:

- [`example_test.go`](https://github.com/gifs/gifs_go/blob/master/examples/import.go)
- [`gifs_test.go`](https://github.com/gifs/gifs-go/blob/master/gifs_test.go)
- [`gifs.go`](https://github.com/gifs/gifs-go/blob/master/gifs.go)

#### Include the package:

```go

package main

import (
	"fmt"

	gifs "github.com/gifs/gifs-go"
)
```

#### Import example

[embedmd]:# (example_test.go go /func ExampleImport.*/ /\n}/)
```go
func ExampleImport() {
	g, err := gifs.New()
	if err != nil {
		fmt.Printf("failed to initialize a new GIFS client, err=%v\n", err)
		return
	}
	param := &gifs.Request{
		URL: "https://www.youtube.com/watch?v=D2EfpQiOQrY",
		Trim: &gifs.Trim{
			Start: 4.5,
			End:   20.5,
		},
		Title: "Migos Dab",
		Tags:  []string{"migos", "dab", "example test"},
	}

	res, err := g.Import(param)
	if err != nil {
		fmt.Printf("failed to import your media, err=%v\n", err)
		return
	}

	if res.Page != "" {
		fmt.Printf("We've got a page alright\n")
	}
	if res.Embed != "" {
		fmt.Printf("We've got an embed page alright\n")
	}
	if res.HasFiles() {
		fmt.Printf("We've got files\n")
	}
	if res.File(gifs.MP4) != "" {
		fmt.Printf("We've got an MP4 file at the bare minimum\n")
	}

	// We've got a page alright
	// We've got an embed page alright
	// We've got files
	// We've got an MP4 file at the bare minimum
}
```

#### ImportBulk example

[embedmd]:# (example_test.go go /func ExampleImportBulk.*/ /\n}/)
```go
func ExampleImportBulk() {
	g, err := gifs.New()
	if err != nil {
		fmt.Printf("failed to initialize a new GIFS client, err=%v\n", err)
		return
	}

	bulkRequest := &gifs.BulkImportRequest{
		ConcurrentImports: 3,
		Requests: []*gifs.Request{
			{
				Title: "Desiigner -- Panda",
				URL:   "https://www.youtube.com/watch?v=E5ONTXHS2mM",
				Tags:  []string{"panda", "broads in atlanta", "phantom", "desiigner"},
			},
			{
				Title: "She writes her own history.",
				URL:   "https://twitter.com/Nike/status/764611634711105537",
				Tags:  []string{"Nike", "She makes her own history", "Running"},
			},
		},
	}

	responses, err := g.ImportBulk(bulkRequest)
	if err != nil {
		fmt.Printf("Failed to bulk import; err=%v\n", err)
		return
	}

	resLen, bulksLen := len(responses), len(bulkRequest.Requests)
	if resLen != bulksLen {
		fmt.Printf("responsesLength(%d) does not match requestsLength (%d)\n", resLen, bulksLen)
		return
	}
}
```
