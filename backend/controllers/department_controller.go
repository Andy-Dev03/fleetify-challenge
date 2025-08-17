package controllers

import (
	"fleetify-backend/config"
	"fleetify-backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Input untuk form department
type DepartmentFormInput struct {
	DepartmentName     string `form:"department_name"`
	MaxClockInTimeStr  string `form:"max_clock_in_time"`
	MaxClockOutTimeStr string `form:"max_clock_out_time"`
}

// Response struct untuk department dan employee
type EmployeeResp struct {
	ID           uint      `json:"id"`
	EmployeeID   string    `json:"employee_id"`
	DepartmentID uint      `json:"department_id"`
	Name         string    `json:"name"`
	Address      string    `json:"address"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
type DepartmentResp struct {
	ID              uint           `json:"id"`
	DepartmentName  string         `json:"department_name"`
	MaxClockInTime  string         `json:"max_clock_in_time"`
	MaxClockOutTime string         `json:"max_clock_out_time"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	Employees       []EmployeeResp `json:"employees"`
}

type DepartmentOnlyResp struct {
	ID              uint      `json:"id"`
	DepartmentName  string    `json:"department_name"`
	MaxClockInTime  string    `json:"max_clock_in_time"`
	MaxClockOutTime string    `json:"max_clock_out_time"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// GetAllDepartments
func GetAllDepartments(c *gin.Context) {
	var departments []models.Department
	if err := config.DB.Preload("Employees").Find(&departments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	var resp []DepartmentResp
	for _, dept := range departments {
		var employees []EmployeeResp
		for _, emp := range dept.Employees {
			employees = append(employees, EmployeeResp{
				ID:           emp.ID,
				EmployeeID:   emp.EmployeeID,
				DepartmentID: emp.DepartmentID,
				Name:         emp.Name,
				Address:      emp.Address,
				CreatedAt:    emp.CreatedAt,
				UpdatedAt:    emp.UpdatedAt,
			})
		}
		resp = append(resp, DepartmentResp{
			ID:              dept.ID,
			DepartmentName:  dept.DepartmentName,
			MaxClockInTime:  dept.MaxClockInTime,
			MaxClockOutTime: dept.MaxClockOutTime,
			CreatedAt:       dept.CreatedAt,
			UpdatedAt:       dept.UpdatedAt,
			Employees:       employees,
		})
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// GetDepartmentDetail
func GetDepartmentDetail(c *gin.Context) {
	id := c.Param("id")
	var department models.Department
	if err := config.DB.Preload("Employees").First(&department, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
		return
	}

	var employees []EmployeeResp
	for _, emp := range department.Employees {
		employees = append(employees, EmployeeResp{
			ID:           emp.ID,
			EmployeeID:   emp.EmployeeID,
			DepartmentID: emp.DepartmentID,
			Name:         emp.Name,
			Address:      emp.Address,
			CreatedAt:    emp.CreatedAt,
			UpdatedAt:    emp.UpdatedAt,
		})
	}

	resp := DepartmentResp{
		ID:              department.ID,
		DepartmentName:  department.DepartmentName,
		MaxClockInTime:  department.MaxClockInTime,
		MaxClockOutTime: department.MaxClockOutTime,
		CreatedAt:       department.CreatedAt,
		UpdatedAt:       department.UpdatedAt,
		Employees:       employees,
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// CreateDepartment
func CreateDepartment(c *gin.Context) {
	var input DepartmentFormInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	// Custom validation
	if input.DepartmentName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Department Name is required"})
		return
	}
	if input.MaxClockInTimeStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Max Clock In Time is required"})
		return
	}
	if input.MaxClockOutTimeStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Max Clock Out Time is required"})
		return
	}

	// Validate time format (HH:mm)
	layout := "15:04"
	clockIn, errIn := time.Parse(layout, input.MaxClockInTimeStr)
	clockOut, errOut := time.Parse(layout, input.MaxClockOutTimeStr)
	if errIn != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid format for max_clock_in_time, expected HH:mm"})
		return
	}
	if errOut != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid format for max_clock_out_time, expected HH:mm"})
		return
	}

	dept := models.Department{
		DepartmentName:  input.DepartmentName,
		MaxClockInTime:  clockIn.Format("15:04:05"),
		MaxClockOutTime: clockOut.Format("15:04:05"),
	}
	if err := config.DB.Create(&dept).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	resp := DepartmentResp{
		ID:              dept.ID,
		DepartmentName:  dept.DepartmentName,
		MaxClockInTime:  dept.MaxClockInTime,
		MaxClockOutTime: dept.MaxClockOutTime,
		CreatedAt:       dept.CreatedAt,
		UpdatedAt:       dept.UpdatedAt,
		Employees:       []EmployeeResp{},
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// Update
func UpdateDepartment(c *gin.Context) {
	id := c.Param("id")
	var department models.Department
	if err := config.DB.First(&department, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
		return
	}

	var input DepartmentFormInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	// Custom validation
	if input.DepartmentName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Department Name is required"})
		return
	}
	if input.MaxClockInTimeStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Max Clock In Time is required"})
		return
	}
	if input.MaxClockOutTimeStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Max Clock Out Time is required"})
		return
	}

	// Validate time format (HH:mm)
	layout := "15:04"
	clockIn, errIn := time.Parse(layout, input.MaxClockInTimeStr)
	clockOut, errOut := time.Parse(layout, input.MaxClockOutTimeStr)
	if errIn != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid format for max_clock_in_time, expected HH:mm"})
		return
	}
	if errOut != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid format for max_clock_out_time, expected HH:mm"})
		return
	}

	// Update fields
	department.DepartmentName = input.DepartmentName
	department.MaxClockInTime = clockIn.Format("15:04:05")
	department.MaxClockOutTime = clockOut.Format("15:04:05")

	if err := config.DB.Save(&department).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	resp := DepartmentOnlyResp{
		ID:              department.ID,
		DepartmentName:  department.DepartmentName,
		MaxClockInTime:  department.MaxClockInTime,
		MaxClockOutTime: department.MaxClockOutTime,
		CreatedAt:       department.CreatedAt,
		UpdatedAt:       department.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// Delete
func DeleteDepartment(c *gin.Context) {
	id := c.Param("id")
	var department models.Department
	if err := config.DB.First(&department, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
		return
	}

	// Hapus semua employee dan attendance terkait
	var employees []models.Employee
	if err := config.DB.Where("department_id = ?", department.ID).Find(&employees).Error; err == nil {
		for _, emp := range employees {
			config.DB.Where("employee_id = ?", emp.EmployeeID).Delete(&models.AttendanceHistory{})
			config.DB.Where("employee_id = ?", emp.EmployeeID).Delete(&models.Attendance{})
			config.DB.Delete(&emp)
		}
	}

	// Hapus department
	if err := config.DB.Delete(&department).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Department deleted successfully"})
}
