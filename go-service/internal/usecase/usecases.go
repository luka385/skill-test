package usecase

import (
	"bytes"
	"log"

	"github.com/luke385/skill-test/internal/ports"
)

type ReportUsecase struct {
	Repo ports.StudentRepository
	Gen  ports.FileGenerator
}

func NewReportUseCase(repo ports.StudentRepository, gen ports.FileGenerator) *ReportUsecase {
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
