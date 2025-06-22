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
	style := sanitizeStyle(req.Style)
	filename := "templates/" + style + ".svg.tmpl"

	tmpl, err := template.ParseFS(templateFS, filename)
	if err != nil {
		log.Println("Template parse error:", err)
		return ""
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

func sanitizeStyle(style string) string {
	_, err := templateFS.Open(fmt.Sprintf("templates/%s.svg", style))
	if err != nil {
		return "default"
	}

	return style
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
