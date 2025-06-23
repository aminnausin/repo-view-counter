package badge

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"repo-view-counter/internal/request"

	"github.com/dustin/go-humanize"
)

// px Constants for calculating width
const charWidth = 6.2
const padding = 10.0

//go:embed templates/*
var templateFS embed.FS

var templates = template.Must(template.ParseFS(templateFS, "templates/*.tmpl"))

type BadgeData struct {
	Label  string
	LabelW string
	LabelX string
	Views  string
	ViewsW string
	ViewsX string
	Colour string
	Width  string
}

func RenderBadgeSVG(req request.BadgeRequest, views int) string {
	filename := req.Style + ".svg.tmpl"

	tmpl := templates.Lookup(filename)
	if tmpl == nil {
		tmpl = templates.Lookup("default.svg.tmpl") // fallback
	}

	viewsStr := humanize.Comma(int64(views))

	totalW, labelW, viewsW := calcSVGWidth(req.Label, viewsStr)

	data := BadgeData{
		Label:  req.Label,
		Views:  viewsStr,
		Colour: req.Colour,
		Width:  fmt.Sprintf("%.1f", totalW),
		LabelW: fmt.Sprintf("%.1f", labelW),
		LabelX: fmt.Sprintf("%.1f", labelW/2),
		ViewsW: fmt.Sprintf("%.1f", viewsW),
		ViewsX: fmt.Sprintf("%.1f", labelW+(viewsW/2)),
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Println("Template execute error:", err)
		return ""
	}

	log.Printf("Got %d views for %s\n", views, req.Repository)

	return buf.String()
}

func estimateTextWidth(text string) (width float64) {
	return float64(len(text)) * charWidth
}

func calcSVGWidth(label, views string) (total float64, labelWidth float64, viewsWidth float64) {
	labelWidth = estimateTextWidth(label) + padding
	viewsWidth = estimateTextWidth(views) + padding
	total = labelWidth + viewsWidth
	return
}
