package model

// PDFRequest payload utama
type PDFRequest struct {
	Config  *PDFConfigDTO `json:"config,omitempty"` // Opsional
	Entries []Entry       `json:"entries"`
}

type Entry struct {
	Name   string `json:"name"`
	Number string `json:"number"`
}

// PDFConfigDTO (Input JSON)
// Menggunakan Pointer supaya bisa membedakan antara "Nilai 0" dan "Tidak Diisi (Nil)"
type PDFConfigDTO struct {
	PageMargin       *float64 `json:"page_margin"`
	NameOverflowMode *string  `json:"name_overflow_mode"`
	MinFontSize      *float64 `json:"min_font_size"`
	MaxNameLines     *int     `json:"max_name_lines"`
	BoxWidth         *float64 `json:"box_width"`
	BoxHeight        *float64 `json:"box_height"`
	InnerOffset      *float64 `json:"inner_offset"`
	OuterBorderWidth *float64 `json:"outer_border_width"`
	InnerBorderOuter *float64 `json:"inner_border_outer"`
	InnerBorderInner *float64 `json:"inner_border_inner"`
	CornerRadius     *float64 `json:"corner_radius"`
	FontName         *string  `json:"font_name"`
	FontSize         *float64 `json:"font_size"`
	LineSpacing      *float64 `json:"line_spacing"`
	NamePrefix       *string  `json:"name_prefix"`
	NumberPrefix     *string  `json:"number_prefix"`
	Columns          *int     `json:"columns"`
}

// PDFConfig (Internal Logic)
// Tanpa pointer agar mudah dikalkulasi di service
type PDFConfig struct {
	PageMargin       float64
	NameOverflowMode string
	MinFontSize      float64
	MaxNameLines     int
	BoxWidth         float64
	BoxHeight        float64
	InnerOffset      float64
	OuterBorderWidth float64
	InnerBorderOuter float64
	InnerBorderInner float64
	CornerRadius     float64
	FontName         string
	FontSize         float64
	LineSpacing      float64
	NamePrefix       string
	NumberPrefix     string
	Columns          int
}
