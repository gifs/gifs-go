package gifs

type MediaType uint

const (
	MP4 MediaType = 1 << iota
	JPG
	GIF
)

func (mt MediaType) Extension() string {
	switch mt {
	default:
		return ""
	case MP4:
		return "mp4"
	case JPG:
		return "jpg"
	case GIF:
		return "gif"
	}
}

func (mt MediaType) String() string {
	return mt.Extension()
}
