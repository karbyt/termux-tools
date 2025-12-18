package handler

import (
	"encoding/json"
	"net/http"

	"go-name/internal/model"
	"go-name/internal/service"
)

type PDFHandler struct {
	Service *service.PDFService
}

func NewPDFHandler(s *service.PDFService) *PDFHandler {
	return &PDFHandler{Service: s}
}

func (h *PDFHandler) Generate(w http.ResponseWriter, r *http.Request) {
	// 1. Validasi Method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Decode JSON Request
	var req model.PDFRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if len(req.Entries) == 0 {
		http.Error(w, "Entries cannot be empty", http.StatusBadRequest)
		return
	}

	// 3. Set Header Response (PENTING AGAR BROWSER TAHU INI PDF)
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=namecards.pdf")

	// 4. Generate PDF langsung ke ResponseWriter
	// ResponseWriter di Go mengimplementasikan interface io.Writer
	err := h.Service.GeneratePDF(w, req)
	if err != nil {
		http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
		return
	}
}
