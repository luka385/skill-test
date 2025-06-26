package report

import (
	"bytes"
	"log"

	"github.com/luke385/skill-test/internal/report/usecase/domain"
)

type StudentRepository interface {
	GetByID(string) (*domain.Student, error)
}

type ReportUsecase struct {
	Repo StudentRepository
	Gen  FileGenerator
}

func NewReportUseCase(repo StudentRepository, gen FileGenerator) *ReportUsecase {
	return &ReportUsecase{Repo: repo, Gen: gen}
}

func (uc *ReportUsecase) Execute(id string) (*bytes.Buffer, string, string, error) {
	student, err := uc.Repo.GetByID(id)
	if err != nil {
		log.Printf("ERROR (GetByID): %v", err)
		return nil, "", "", err
	}

	buf, err := uc.Gen.Generate(student)
	if err != nil {
		log.Printf("ERROR (Generate): %v", err)
		return nil, "", "", err
	}

	contentType := uc.Gen.GetContentType()
	ext := uc.Gen.GetFileExtension()

	return buf, contentType, ext, nil
}
