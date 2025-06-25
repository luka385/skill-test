package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/luke385/skill-test/internal/usecase"
)

type StudentHandler struct {
	usecase *usecase.ReportUsecase
}

func NewStudentHandler(uc *usecase.ReportUsecase) *StudentHandler {
	return &StudentHandler{uc}
}

func (h *StudentHandler) GenerateReport(c *gin.Context) {
	id := c.Param("id")
	log.Printf("GenerateReport called with id=%s", id)
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Student ID is required"})
		return
	}

	pdfBuf, err := h.usecase.Execute(id)
	if err != nil {
		log.Printf("ERROR (GenerateReport): %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=student_"+id+".pdf")
	c.Data(http.StatusOK, "application/pdf", pdfBuf.Bytes())
}
