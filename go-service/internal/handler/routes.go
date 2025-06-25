package handler

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, h *StudentHandler) {
	r.GET("/api/v1/students/:id/report", h.GenerateReport)
}
