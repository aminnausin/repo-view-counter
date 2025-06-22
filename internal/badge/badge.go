package badge

import (
	"fmt"
	"net/http"
	"regexp"
	"repo-view-counter/internal/db"
	"repo-view-counter/internal/request"
)

var repoPattern = regexp.MustCompile(`^[a-zA-Z0-9-]+/[a-zA-Z0-9._-]+$`)

type Service struct {
	db db.Database
}

func NewService(db db.Database) Service {
	return Service{db: db}
}

func (s Service) HandleBadge(req request.BadgeRequest) (string, error) {
	if req.Repository == "" || !repoPattern.MatchString(req.Repository) {
		return "", fmt.Errorf("no repo provided")
	}

	if isGitHubCamo(req.UserAgent) {
		_ = s.db.IncrementViews(req.Repository)
	}

	views, err := s.db.GetViews(req.Repository)
	if err != nil {
		return "", err
	}

	svg := RenderBadgeSVG(req, views)
	return svg, nil
}

func Handler(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := request.NewBadgeRequest(r)

		if req.Repository == "" || !repoPattern.MatchString(req.Repository) {
			http.Redirect(w, r, "https://github.com/aminnausin/repo-view-counter", http.StatusPermanentRedirect)
			return
		}

		svg, err := s.HandleBadge(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/svg+xml")
		w.Header().Set("Cache-Control", "no-store, must-revalidate")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(svg))
	}
}

func isGitHubCamo(agent string) bool {
	return len(agent) >= 11 && agent[:11] == "github-camo"
}
