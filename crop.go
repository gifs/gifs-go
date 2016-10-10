package gifs

type Crop struct {
	X float32 `json:"x,omitempty"`
	Y float32 `json:"y,omitempty"`

	Height float32 `json:"height,omitempty"`
	Width  float32 `json:"width,omitempty"`
}
