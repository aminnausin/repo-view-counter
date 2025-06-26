package badge

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"repo-view-counter/internal/db"
	"repo-view-counter/internal/request"
	"strconv"
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
		w.Header().Set("Cache-Control", "public, max-age=5, stale-while-revalidate=5")
		w.Header().Set("Content-Length", strconv.Itoa(len(svg)))
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(svg))

		if err != nil {
			log.Printf("Error returning image: %s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func isGitHubCamo(agent string) bool {
	return len(agent) >= 11 && agent[:11] == "github-camo"
}
