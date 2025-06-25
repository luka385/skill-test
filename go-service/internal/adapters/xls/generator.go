package xls

import (
	"bytes"

	"github.com/xuri/excelize/v2"

	"github.com/luke385/skill-test/internal/domain"
	"github.com/luke385/skill-test/internal/ports"
)

type ExcelAdapter struct{}

func NewXLSGenerator() ports.FileGenerator {
	return &ExcelAdapter{}
}

func (x *ExcelAdapter) Generate(s *domain.Student) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	sheet := "Sheet1"

	headers := []string{"Name", "Email", "Class"}
	for colIdx, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(colIdx+1, 1)
		if err := f.SetCellValue(sheet, cell, h); err != nil {
			return nil, err
		}
	}

	values := []interface{}{s.Name, s.Email, s.Class}
	for colIdx, v := range values {
		cell, _ := excelize.CoordinatesToCellName(colIdx+1, 2)
		if err := f.SetCellValue(sheet, cell, v); err != nil {
			return nil, err
		}
	}

	buf := new(bytes.Buffer)
	if err := f.Write(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

func (x *ExcelAdapter) GetContentType() string {
	return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
}

func (x *ExcelAdapter) GetFileExtension() string {
	return "xlsx"
}
