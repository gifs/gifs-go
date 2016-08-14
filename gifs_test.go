package gifs_go_test

import (
	"testing"

	gifs "github.com/gifs/gifs_go"
)

func newClient(t *testing.T) *gifs.GIFS {
	g, err := gifs.New()
	if err != nil {
		t.Fatal(err)
	}

	return g
}

func TestNew(t *testing.T) {
	g := newClient(t)
	if g == nil {
		t.Errorf("expected non-nil GIFS handler")
	}
}

func TestImportSources(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip()
	}

	g := newClient(t)
	sources := []string{
		"https://twitter.com/kanyewest/status/726835785274646529",
		"https://www.youtube.com/watch?v=jNQXAC9IVRw",
		"https://www.facebook.com/Pagefanclub/videos/1070860819621603/",
		"https://j.gifs.com/PNoDGy.gif",
		"http://video.webmfiles.org/elephants-dream.webm",
		"http://www.engr.colostate.edu/me/facil/dynamics/files/flame.avi",
		"https://github.com/commaai/research/raw/master/images/drive_simulator.gif",
		"https://www.youtube.com/watch?v=_gB2iWln0ls",
	}

	responses, err := g.ImportSources(sources...)
	if err != nil {
		t.Errorf("%v", err)
	}

	if len(responses) < 1 {
		t.Fatalf("did not get back any responses")
	}

	if respLen, srcLen := len(responses), len(sources); respLen != srcLen {
		t.Fatalf("responseLen=%d wanted %d", respLen, srcLen)
	}

	for i, resp := range responses {
		if resp == nil {
			t.Errorf("#%d: resp is nil;source: %s", i, sources[i])
			continue
		}
		if resp.Page == "" {
			t.Errorf("#%d: attribute Page is empty;source: %s", i, sources[i])
		}
	}
}

func TestImport(t *testing.T) {
	t.Parallel()

	param := gifs.Params{
		URL:   "https://camo.githubusercontent.com/cd1f1a4b10bb14133ae48db167919c418d455537/68747470733a2f2f73746f726167652e676f6f676c65617069732e636f6d2f63646e2e676966732e636f6d2f67656e69652d6769746875622d616e696d6174696f6e2e6769663f763d34",
		Title: "GIFS Genie",
		Tags:  []string{"tests", "hello golang"},
	}
	g := newClient(t)
	resp, err := g.Import(&param)
	if err != nil {
		t.Errorf("err=%v want nil error", err)
	}
	t.Logf("resp=%v err=%v\n", resp, err)
	if resp == nil {
		t.Errorf("expected a non-nil response")
		return
	}
	if resp.Page == "" {
		t.Errorf("expected a non-empty page")
	}
	if resp.OEmbed == "" {
		t.Errorf("expected a non-blank OEmbed URL")
	}
	if resp.Embed == "" {
		t.Errorf("expected a non-blank OEmbed URL")
	}

	if resp.File(gifs.MP4) == "" {
		t.Errorf("expected at least an MP4 back")
	}
}
