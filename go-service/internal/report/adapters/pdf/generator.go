package pdf

import (
	"bytes"

	"github.com/jung-kurt/gofpdf"

	"github.com/luke385/skill-test/internal/report"
	"github.com/luke385/skill-test/internal/report/usecase/domain"
)

type PDFAdapter struct{}

func NewPDFAdapter() report.FileGenerator {
	return &PDFAdapter{}
}

func (p *PDFAdapter) Generate(s *domain.Student) (*bytes.Buffer, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "Student Report", "", 1, "C", false, 0, "")
	pdf.Ln(6)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 8, "Name:")
	pdf.Cell(0, 8, s.Name)
	pdf.Ln(8)

	pdf.Cell(40, 8, "Email:")
	pdf.Cell(0, 8, s.Email)
	pdf.Ln(8)

	pdf.Cell(40, 8, "Class:")
	pdf.Cell(0, 8, s.Class)
	pdf.Ln(8)

	buf := new(bytes.Buffer)
	if err := pdf.Output(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

func (p *PDFAdapter) GetContentType() string {
	return "application/pdf"
}

func (p *PDFAdapter) GetFileExtension() string {
	return "pdf"
}
