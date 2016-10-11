package gifs

// Crop holds the offsets and final dimensions of the
// desired media after cropping. Both X and Y offsets 
// will tell us where the top-left corner of the cropped
// media will be. Then Height and Width can determine the 
// top-right, bottom-left and bottom-right coordinates. 
type Crop struct {
	// X is the horizontal axis offset from the left
	X float32 `json:"x,omitempty"`
	// Y is the vertical axis offset from the top 
	Y float32 `json:"y,omitempty"`
	// Height of the desired media after cropping
	Height float32 `json:"height,omitempty"`
	// Width of the desired media after cropping
	Width  float32 `json:"width,omitempty"`
}
