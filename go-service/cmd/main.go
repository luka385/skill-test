package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/luke385/skill-test/internal/adapters/client"
	"github.com/luke385/skill-test/internal/adapters/pdf"
	"github.com/luke385/skill-test/internal/handler"
	"github.com/luke385/skill-test/internal/usecase"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using OS environment vars")
	}
	log.Printf("NODE_API_URL=%s NODE_API_USER=%s NODE_API_PASS=%s CSRF_COOKIE_NAME=%s", os.Getenv("NODE_API_URL"), os.Getenv("NODE_API_USER"), os.Getenv("NODE_API_PASS"), os.Getenv("CSRF_COOKIE_NAME"))

	// Init client
	repo, err := client.NewNodeAPIClient()
	if err != nil {
		log.Fatalf("Error initializing NodeAPIClient: %v", err)
	}

	// Dependencies
	pdfGen := pdf.NewPDFAdapter()
	uc := usecase.NewReportUseCase(repo, pdfGen)
	h := handler.NewStudentHandler(uc)

	// Gin
	r := gin.Default()

	// Routes
	handler.RegisterRoutes(r, h)
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Not found"})
	})

	// Run server
	log.Println("Go PDF report microservice running at http://localhost:8080")
	log.Println("Try GET /api/v1/students/:id/report")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
