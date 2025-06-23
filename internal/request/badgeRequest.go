package request

import (
	"net/http"
	"strings"
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

	label := strings.ToLower(strings.TrimSpace(query.Get("label")))
	if label == "" {
		label = "views"
	}

	style := strings.ToLower(query.Get("style"))
	if style == "" {
		style = "default"
	}

	colour := "#" + query.Get("colour")
	if query.Get("colour") == "" {
		colour = "#007ec6"
	}

	return BadgeRequest{
		UserAgent:  r.UserAgent(),
		Repository: strings.ToLower(strings.TrimSpace(query.Get("repo"))),
		Label:      label,
		Colour:     colour,
		Style:      style,
	}
}
