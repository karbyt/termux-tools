package service

import (
	"fmt"
	"io"
	"strings"

	"go-name/internal/model"

	"github.com/go-pdf/fpdf"
)

// Default Configuration
var DefaultConfig = model.PDFConfig{
	PageMargin:       5.0,
	NameOverflowMode: "AUTO_SHRINK",
	MinFontSize:      8.0,
	MaxNameLines:     2,
	BoxWidth:         100.0,
	BoxHeight:        38.0,
	InnerOffset:      3.0,
	OuterBorderWidth: 0.1,
	InnerBorderOuter: 0.7,
	InnerBorderInner: 0.3,
	CornerRadius:     3.0,
	FontName:         "Times",
	FontSize:         20.0,
	LineSpacing:      1.0,
	NamePrefix:       "",
	NumberPrefix:     "",
	Columns:          2,
}

type PDFService struct{}

func NewPDFService() *PDFService {
	return &PDFService{}
}

// GeneratePDF menulis binary PDF ke writer (io.Writer)
func (s *PDFService) GeneratePDF(w io.Writer, req model.PDFRequest) error {
	// 1. Merge Config: Start with Default, Override with User Input
	cfg := s.mergeConfig(req.Config)

	// 2. Setup PDF dengan Config yang sudah final
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(cfg.PageMargin, cfg.PageMargin, cfg.PageMargin)
	pdf.SetAutoPageBreak(false, cfg.PageMargin)
	pdf.AddPage()

	_, pageHeight := pdf.GetPageSize()

	xStart := cfg.PageMargin
	yStart := cfg.PageMargin
	x, y := xStart, yStart
	col := 0

	for _, entry := range req.Entries {
		if col == cfg.Columns {
			col = 0
			x = xStart
			y += cfg.BoxHeight
		}

		// Cek ganti halaman
		if y+cfg.BoxHeight > pageHeight-cfg.PageMargin {
			pdf.AddPage()
			y = yStart
			x = xStart
			col = 0
		}

		s.drawBox(pdf, cfg, x, y, entry.Name, entry.Number)
		x += cfg.BoxWidth
		col++
	}

	return pdf.Output(w)
}

func (s *PDFService) mergeConfig(dto *model.PDFConfigDTO) model.PDFConfig {
	cfg := DefaultConfig // Copy by value

	if dto == nil {
		return cfg
	}

	// Cek satu per satu field, jika tidak nil, timpa.
	if dto.PageMargin != nil {
		cfg.PageMargin = *dto.PageMargin
	}
	if dto.NameOverflowMode != nil {
		cfg.NameOverflowMode = *dto.NameOverflowMode
	}
	if dto.MinFontSize != nil {
		cfg.MinFontSize = *dto.MinFontSize
	}
	if dto.MaxNameLines != nil {
		cfg.MaxNameLines = *dto.MaxNameLines
	}
	if dto.BoxWidth != nil {
		cfg.BoxWidth = *dto.BoxWidth
	}
	if dto.BoxHeight != nil {
		cfg.BoxHeight = *dto.BoxHeight
	}
	if dto.InnerOffset != nil {
		cfg.InnerOffset = *dto.InnerOffset
	}
	if dto.OuterBorderWidth != nil {
		cfg.OuterBorderWidth = *dto.OuterBorderWidth
	}
	if dto.InnerBorderOuter != nil {
		cfg.InnerBorderOuter = *dto.InnerBorderOuter
	}
	if dto.InnerBorderInner != nil {
		cfg.InnerBorderInner = *dto.InnerBorderInner
	}
	if dto.CornerRadius != nil {
		cfg.CornerRadius = *dto.CornerRadius
	}
	if dto.FontName != nil {
		cfg.FontName = *dto.FontName
	}
	if dto.FontSize != nil {
		cfg.FontSize = *dto.FontSize
	}
	if dto.LineSpacing != nil {
		cfg.LineSpacing = *dto.LineSpacing
	}
	if dto.NamePrefix != nil {
		cfg.NamePrefix = *dto.NamePrefix
	}
	if dto.NumberPrefix != nil {
		cfg.NumberPrefix = *dto.NumberPrefix
	}
	if dto.Columns != nil {
		cfg.Columns = *dto.Columns
	}

	return cfg
}

func (s *PDFService) drawBox(pdf *fpdf.Fpdf, cfg model.PDFConfig, x, y float64, name, number string) {
	// 1. Border Luar
	pdf.SetLineWidth(cfg.OuterBorderWidth)
	pdf.Rect(x, y, cfg.BoxWidth, cfg.BoxHeight, "D")

	// 2. Border Dalam
	innerX := x + cfg.InnerOffset
	innerY := y + cfg.InnerOffset
	innerW := cfg.BoxWidth - (2 * cfg.InnerOffset)
	innerH := cfg.BoxHeight - (2 * cfg.InnerOffset)

	pdf.SetLineWidth(cfg.InnerBorderOuter)
	pdf.RoundedRect(innerX, innerY, innerW, innerH, cfg.CornerRadius, "1234", "D")

	gap := 1.5
	pdf.SetLineWidth(cfg.InnerBorderInner)
	pdf.RoundedRect(innerX+gap, innerY+gap, innerW-(2*gap), innerH-(2*gap), cfg.CornerRadius-1, "1234", "D")

	// 3. Text Logic
	availableWidth := innerW - 4.0
	nameText := cfg.NamePrefix + name
	numberText := cfg.NumberPrefix + number

	currentFontSize := cfg.FontSize
	var nameLines []string

	// Overflow Handler
	switch cfg.NameOverflowMode {
	case "ABBREVIATE":
		pdf.SetFont(cfg.FontName, "", currentFontSize)
		if pdf.GetStringWidth(nameText) > availableWidth {
			nameText = cfg.NamePrefix + abbreviateName(name)
		}
		nameLines = []string{nameText}
	case "WRAP_LINE":
		nameLines = wrapText(pdf, nameText, availableWidth, cfg.FontName, cfg.FontSize, cfg.MaxNameLines)
	case "AUTO_SHRINK":
		pdf.SetFont(cfg.FontName, "", currentFontSize)
		for pdf.GetStringWidth(nameText) > availableWidth && currentFontSize > cfg.MinFontSize {
			currentFontSize -= 0.5
			pdf.SetFont(cfg.FontName, "", currentFontSize)
		}
		nameLines = []string{nameText}
	default:
		nameLines = []string{nameText}
	}

	totalLines := append(nameLines, numberText)

	// 4. Center Calculation
	mmPerPt := 0.352778
	totalBlockHeight := 0.0
	var lineHeights []float64

	for i := range totalLines {
		fs := currentFontSize
		if i >= len(nameLines) {
			fs = cfg.FontSize
		}
		lh := (fs * mmPerPt) * cfg.LineSpacing
		lineHeights = append(lineHeights, lh)
		totalBlockHeight += lh
	}

	cursorY := y + (cfg.BoxHeight-totalBlockHeight)/2

	for i, line := range totalLines {
		fs := currentFontSize
		if i >= len(nameLines) {
			fs = cfg.FontSize
		}
		pdf.SetFont(cfg.FontName, "", fs)
		pdf.SetXY(x, cursorY)
		pdf.CellFormat(cfg.BoxWidth, lineHeights[i], line, "", 0, "C", false, 0, "")
		cursorY += lineHeights[i]
	}
}

// Helper functions (Private)
func abbreviateName(name string) string {
	parts := strings.Fields(name)
	if len(parts) <= 2 {
		return name
	}
	first, second := parts[0], parts[1]
	var rest []string
	for _, p := range parts[2:] {
		if len(p) > 0 {
			rest = append(rest, string(p[0])+".")
		}
	}
	return fmt.Sprintf("%s %s %s", first, second, strings.Join(rest, " "))
}

func wrapText(pdf *fpdf.Fpdf, text string, maxWidth float64, fontName string, fontSize float64, maxLines int) []string {
	words := strings.Fields(text)
	var lines []string
	currentLine := ""
	pdf.SetFont(fontName, "", fontSize)
	for _, word := range words {
		testLine := word
		if currentLine != "" {
			testLine = currentLine + " " + word
		}
		if pdf.GetStringWidth(testLine) <= maxWidth {
			currentLine = testLine
		} else {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = word
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}
	if len(lines) > maxLines {
		return lines[:maxLines]
	}
	return lines
}
