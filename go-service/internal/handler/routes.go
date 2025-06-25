package handler

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, h *StudentHandler) {
	api := r.Group("/api/v1")
	{
		api.GET("/students/:id/report", h.GenerateReport)
	}
}
