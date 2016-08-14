package gifs_go_test

import (
	"fmt"

	gifs "github.com/gifs/gifs_go"
)

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

func ExampleImportBulk() {
	g, err := gifs.New()
	if err != nil {
		fmt.Printf("failed to initialize a new GIFS client, err=%v\n", err)
		return
	}

	bulkParams := &gifs.BulkImportParams{
		ConcurrentImports: 3,
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
		fmt.Printf("Failed to bulk import; err=%v\n", err)
		return
	}

	resLen, srcsLen := len(responses), len(sources)
	if resLen != srcsLen {
		fmt.Printf("responsesLength(%d) does not match requestsLength (%d)\n", resLen, srcsLen)
		return
	}
}
