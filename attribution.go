package gifs

// Attribution credits creators with importing media
// e.g the user on GIFS.com that first identified BigFoot
type Attribution struct {
	// SiteName is the title of the site that
	// this attribution  is made from.
	SiteName string `json:"site,omitempty"`
	// SiteURL is the origin of the site that this attribution
	// is being made from.
	SiteURL string `json:"url,omitempty"`
	// SiteUsername is the unique identifier for the author/user
	// that data is being attributed/credited to.
	SiteUsername string `json:"user,omitempty"`
}
