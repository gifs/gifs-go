## gifs\_go

[![GoDoc](https://godoc.org/github.com/gifs/gifs_go?status.svg)](https://godoc.org/github.com/gifs/gifs_go)

Go package for interacting with the gifs.com API for transcoding and importing
media to make at least .MP4, .JPG, .GIF outputs.

Import all your videos by uploading and by passing in any source.

It is super easy to import media and integrate in your code.
Please see the examples below on how to import a single source, then in bulk
each in one call. You can also see more usage and convenience methods in files:
- `examples_test.go`
- `gifs_test.go`
- `gifs.go`

### Examples

For the examples below, include these imports as the predicates to have runnable examples:
```go

package main

import (
	"fmt"

	gifs "github.com/gifs/gifs_go"
)
```

##### Show me the code:

##### Import example
```go
func ExampleImport() {
	g, err := gifs.New()
	if err != nil {
		fmt.Printf("failed to initialize a new GIFS client, err=%v\n", err)
		return
	}
	param := &gifs.Params{
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

##### ImportBulk example
```go
func ExampleImportBulk() {
	g, err := gifs.New()
	if err != nil {
		fmt.Printf("failed to initialize a new GIFS client, err=%v\n", err)
		return
	}

	bulkParams := &gifs.BulkImportParams{
		Params: []*gifs.Params{
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

	responses, err := g.ImportBulk(bulkParams)
	if err != nil {
		fmt.Printf("Failed to bulk import; err=%v\n", err)
		return
	}

	resLen, bulksLen := len(responses), len(bulkParams.Params)
	if resLen != bulksLen {
		fmt.Printf("responsesLength(%d) does not match requestsLength (%d)\n", resLen, bulksLen)
		return
	}
}
```
