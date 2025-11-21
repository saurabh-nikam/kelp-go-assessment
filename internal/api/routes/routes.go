package routes

import (
	"kelp-go-assessment/internal/api/handlers"
	"kelp-go-assessment/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	svc := service.NewDataService()
	h := handlers.NewHandler(svc)

	api := r.Group("/api")
	{
		api.GET("/financials", h.GetFinancials)
		api.GET("/sales", h.GetSales)
		api.GET("/employee", h.GetEmployeeStats)
	}

	return r
}
