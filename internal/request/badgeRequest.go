package request

import (
	"net/http"
)

type BadgeRequest struct {
	UserAgent  string
	Repository string
	Label      string
	Colour     string
	Style      string
}

func NewBadgeRequest(r *http.Request) BadgeRequest {
	query := r.URL.Query()

	label := query.Get("label")
	if label == "" {
		label = "Views"
	}

	style := query.Get("style")
	if style == "" {
		style = "default"
	}

	colour := "#" + query.Get("colour")
	if query.Get("colour") == "" {
		colour = "#007ec6"
	}

	return BadgeRequest{
		UserAgent:  r.UserAgent(),
		Repository: query.Get("repo"),
		Label:      label,
		Colour:     colour,
		Style:      style,
	}
}
