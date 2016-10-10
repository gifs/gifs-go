package gifs_test

import (
	"fmt"
	"log"

	"github.com/gifs/gifs-go"
)

func ExampleImport() {
	g, err := gifs.New()
	if err != nil {
		log.Fatalf("failed to initialize a new GIFS client, err=%v\n", err)
	}

	param := &gifs.Request{
		URL: "https://www.youtube.com/watch?v=Vhh_GeBPOhs",
		Trim: &gifs.Trim{
			Start: 4.5,
			End:   19.5,
		},
		Crop: &gifs.Crop{
			X:      40,
			Y:      10,
			Width:  200,
			Height: 200,
		},

		Title: "Developers developers developers",
		Tags:  []string{"steve balmer", "steve turnt", "developers developers"},

		CreatedFrom: "gifs-go-tests",

		Attribution: &gifs.Attribution{
			SiteName:     "gifs-developers",
			SiteURL:      "https://github.com/gifs",
			SiteUsername: "gifs",
		},
	}

	res, err := g.Import(param)
	if err != nil {
		log.Fatalf("failed to import your media, err=%v\n", err)
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

	// Output:
	// We've got a page alright.
	// We've got an embed page alright
	// We've got files
	// We've got an MP4 file at the bare minimum
}

func ExampleImportBulk() {
	g, err := gifs.New()
	if err != nil {
		log.Fatalf("failed to initialize a new GIFS client, err=%v\n", err)
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
		log.Fatalf("Failed to bulk import; err=%v\n", err)
	}

	resLen, bulksLen := len(responses), len(bulkRequest.Requests)
	if resLen != bulksLen {
		log.Fatalf("responsesLength(%d) does not match requestsLength (%d)\n", resLen, bulksLen)
	}
}

func ExampleImportBySources() {
	g, err := gifs.New()
	if err != nil {
		fmt.Printf("failed to initialize a new GIFS client, err=%v\n", err)
		return
	}

	sources := []string{
		"https://twitter.com/kanyewest/status/726835785274646529",
		"http://www.engr.colostate.edu/me/facil/dynamics/files/flame.avi",
		"https://www.youtube.com/watch?v=jNQXAC9IVRw",
		"https://www.facebook.com/Pagefanclub/videos/1070860819621603/",
		"https://j.gifs.com/PNoDGy.gif",
		"http://video.webmfiles.org/elephants-dream.webm",
		"https://github.com/commaai/research/raw/master/images/drive_simulator.gif",
		"https://www.youtube.com/watch?v=_gB2iWln0ls",
	}

	responses, err := g.ImportSources(sources...)
	if err != nil {
		log.Fatalf("Failed to bulk import; err=%v\n", err)
	}

	resLen, srcsLen := len(responses), len(sources)
	if resLen != srcsLen {
		log.Fatalf("responsesLength(%d) does not match requestsLength (%d)\n", resLen, srcsLen)
	}
}
