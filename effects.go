package gifs

// Effects define the collection of alterations
// that will be applied to media.
// Any timed effect's start and end times should
// be considered relative to time 0 after trimming.
type Effects struct {
	Overlay []*Overlay `json:"overlay,omitempty"`
	Pad     []*Pad     `json:"pad,omitempty"`
	Flip    []*Flip    `json:"flip,omitempty"`
	Invert  []*Invert  `json:"invert,omitempty"`
}

// Timeline defines the duration on a timescale
// for which an effect will appear on the video.
type Timeline struct {
	Start float32 `json:"start,omitempty"`
	End   float32 `json:"end,omitempty"`
}

// Section defines a bounding box in which
// the defined effect will appear.
type Section struct {
	X      string  `json:"x,omitempty"`
	Y      string  `json:"y,omitempty"`
	Width  float32 `json:"width,omitempty"`
	Height float32 `json:"height,omitempty"`
}

// Overlay defines a gif, or static image that will be
// applied to the final media at a position for a
// specified time period if defined.
type Overlay struct {
	X        string    `json:"x,omitempty"`
	Y        string    `json:"y,omitempty"`
	Timeline *Timeline `json:"timeline,omitempty"`
	Source   string    `json:"source,omitempty"`

	// LoopCount if set defines the number of times
	// that an animated overlay will loop for.
	LoopCount int `json:"loop_count,omitempty"`
}

type Pad struct {
	X      float32 `json:"x,omitempty"`
	Y      float32 `json:"y,omitempty"`
	Color  string  `json:"color,omitempty"`
	Height float32 `json:"height,omitempty"`
	Width  float32 `json:"width,omitempty"`
}

type Flip struct {
	Horizontal bool `json:"horizontal,omitempty"`
	Vertical   bool `json:"vertical,omitempty"`
}

type Invert struct {
	Enable   bool      `json:"enable,omitempty"`
	Section  *Section  `json:"section,omitempty"`
	Timeline *Timeline `json:"timeline,omitempty"`
}
