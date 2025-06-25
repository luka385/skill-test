package handler

import (
	"bytes"

	"github.com/gin-gonic/gin"

	"github.com/luke385/skill-test/internal/usecase"
)

type StudentHandler struct {
	pdfUC  *usecase.ReportUsecase
	xlsxUC *usecase.ReportUsecase
}

func NewStudentHandler(pdfUC, xlsxUC *usecase.ReportUsecase) *StudentHandler {
	return &StudentHandler{pdfUC: pdfUC, xlsxUC: xlsxUC}
}

func (h *StudentHandler) GenerateReport(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "Student ID required"})
		return
	}

	format := c.DefaultQuery("format", "pdf")

	var (
		buf         *bytes.Buffer
		contentType string
		ext         string
		err         error
	)
	switch format {
	case "xlsx":
		buf, contentType, ext, err = h.xlsxUC.Execute(id)
	default: // "pdf"
		buf, contentType, ext, err = h.pdfUC.Execute(id)
	}
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "attachment; filename=student_report."+ext)
	c.Data(200, contentType, buf.Bytes())
}
