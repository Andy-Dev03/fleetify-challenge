package routes

import (
	"fleetify-backend/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes untuk registrasi semua route
func RegisterRoutes(app *gin.Engine) {
	api := app.Group("/api") // semua route pakai prefix /api

	// Employee routes
	api.GET("/employees", controllers.GetAllEmployees)
	api.GET("/employee/:id", controllers.GetEmployeeDetail)
	api.POST("/employee", controllers.CreateEmployee)
	api.PATCH("/employee/:id", controllers.UpdateEmployee)
	api.DELETE("/employee/:id", controllers.DeleteEmployee)

	// Departement routes
	api.GET("/departements", controllers.GetAllDepartments)
	api.GET("/departement/:id", controllers.GetDepartmentDetail)
	api.POST("/departement", controllers.CreateDepartment)
	api.PATCH("/departement/:id", controllers.UpdateDepartment)
	api.DELETE("/departement/:id", controllers.DeleteDepartment)

	// Attendance routes
	api.POST("/attendance", controllers.CreateAttendance)
	api.PUT("/attendance/:id", controllers.UpdateAttendance)
	api.GET("/attendance/logs", controllers.GetAttendanceLogs)
}
