package usecase

import (
	"bytes"
	"log"

	"github.com/luke385/skill-test/internal/ports"
)

type ReportUsecase struct {
	Repo ports.StudentRepository
	PDF  ports.PDFGenerator
}

func NewReportUseCase(repo ports.StudentRepository, pdf ports.PDFGenerator) *ReportUsecase {
	return &ReportUsecase{Repo: repo, PDF: pdf}
}

func (uc *ReportUsecase) Execute(id string) (*bytes.Buffer, error) {
	student, err := uc.Repo.GetByID(id)
	if err != nil {
		log.Printf("ERROR (GetByID): %v", err)
		return nil, err
	}
	return uc.PDF.Generate(student)
}
