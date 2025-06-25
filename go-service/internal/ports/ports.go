package ports

import (
	"bytes"

	"github.com/luke385/skill-test/internal/domain"
)

type StudentRepository interface {
	GetByID(string) (*domain.Student, error)
}

type PDFGenerator interface {
	Generate(*domain.Student) (*bytes.Buffer, error)
}
