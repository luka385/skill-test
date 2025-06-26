package report

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/luke385/skill-test/internal/report/usecase/domain"
)

type FileGenerator interface {
	Generate(*domain.Student) (*bytes.Buffer, error)
	GetContentType() string
	GetFileExtension() string
}

type StudentHandler struct {
	pdfUC  *ReportUsecase
	xlsxUC *ReportUsecase
}

func NewStudentHandler(pdfUC, xlsxUC *ReportUsecase) *StudentHandler {
	return &StudentHandler{pdfUC: pdfUC, xlsxUC: xlsxUC}
}

func (h *StudentHandler) GenerateReport(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Student ID required"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filename := "student_" + id + "." + ext
	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, contentType, buf.Bytes())
}

func RegisterRoutes(r *gin.Engine, h *StudentHandler) {
	api := r.Group("/api/v1")
	{
		api.GET("/students/:id/report", h.GenerateReport)
	}
}
