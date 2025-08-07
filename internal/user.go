package internal

import (
	"github.com/mileusna/useragent"
)

type UserAgent struct {
	Browser string
	OS      string
	Version string
	Device  string
}

func ParseUserAgent(userAgent string) *UserAgent {
	ua := useragent.Parse(userAgent)

	return &UserAgent{
		Browser: ua.Name,
		OS:      ua.OS,
		Version: ua.OSVersion,
		Device:  ua.Device,
	}
}
