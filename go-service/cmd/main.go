package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/luke385/skill-test/internal/report"
	"github.com/luke385/skill-test/internal/report/adapters/client"
	"github.com/luke385/skill-test/internal/report/adapters/pdf"
	"github.com/luke385/skill-test/internal/report/adapters/xls"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using OS environment vars")
	}
	log.Printf(
		"NODE_API_URL=%s NODE_API_USER=%s NODE_API_PASS=%s CSRF_COOKIE_NAME=%s",
		os.Getenv("NODE_API_URL"),
		os.Getenv("NODE_API_USER"),
		os.Getenv("NODE_API_PASS"),
		os.Getenv("CSRF_COOKIE_NAME"),
	)

	// Initialize repository client
	repo, err := client.NewNodeAPIClient()
	if err != nil {
		log.Fatalf("Error initializing NodeAPIClient: %v", err)
	}

	// Create UseCases for PDF and Excel
	pdfUC := report.NewReportUseCase(repo, pdf.NewPDFAdapter())
	excelUC := report.NewReportUseCase(repo, xls.NewXLSGenerator())

	// Initialize handler with both usecases
	h := report.NewStudentHandler(pdfUC, excelUC)

	// Set up Gin
	r := gin.Default()

	// Register routes
	report.RegisterRoutes(r, h)
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Not found"})
	})

	// Start server
	log.Println("Report microservice running at http://localhost:8080")
	log.Println("Try GET /api/v1/students/:id/report with Accept: application/pdf or application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
