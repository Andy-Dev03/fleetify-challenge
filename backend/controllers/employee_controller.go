package controllers

import (
	"fleetify-backend/config"
	"fleetify-backend/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response structs
type EmployeeDepartmentResp struct {
	ID              uint   `json:"id"`
	DepartmentName  string `json:"department_name"`
	MaxClockInTime  string `json:"max_clock_in_time"`
	MaxClockOutTime string `json:"max_clock_out_time"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type EmployeeDetailResp struct {
	ID           uint                   `json:"id"`
	EmployeeID   string                 `json:"employee_id"`
	DepartmentID uint                   `json:"department_id"`
	Name         string                 `json:"name"`
	Address      string                 `json:"address"`
	CreatedAt    string                 `json:"created_at"`
	UpdatedAt    string                 `json:"updated_at"`
	Department   EmployeeDepartmentResp `json:"department"`
}

// Helper to convert Department model to response
func toEmployeeDepartmentResp(dept models.Department) EmployeeDepartmentResp {
	return EmployeeDepartmentResp{
		ID:              dept.ID,
		DepartmentName:  dept.DepartmentName,
		MaxClockInTime:  dept.MaxClockInTime,
		MaxClockOutTime: dept.MaxClockOutTime,
		CreatedAt:       dept.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       dept.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toEmployeeDetailResp(emp models.Employee) EmployeeDetailResp {
	return EmployeeDetailResp{
		ID:           emp.ID,
		EmployeeID:   emp.EmployeeID,
		DepartmentID: emp.DepartmentID,
		Name:         emp.Name,
		Address:      emp.Address,
		CreatedAt:    emp.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    emp.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Department:   toEmployeeDepartmentResp(emp.Department),
	}
}

// Get all employees
func GetAllEmployees(c *gin.Context) {
	var employees []models.Employee
	if err := config.DB.Preload("Department").Find(&employees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	var resp []EmployeeDetailResp
	for _, emp := range employees {
		resp = append(resp, toEmployeeDetailResp(emp))
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// Get employee detail by ID
func GetEmployeeDetail(c *gin.Context) {
	id := c.Param("id")
	var employee models.Employee
	if err := config.DB.Preload("Department").First(&employee, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": toEmployeeDetailResp(employee)})
}

// Create a new employee
func CreateEmployee(c *gin.Context) {
	var input struct {
		DepartmentID uint   `form:"department_id"`
		Name         string `form:"name"`
		Address      string `form:"address"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	// Custom validation
	if input.DepartmentID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Department is required"})
		return
	}
	if input.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}
	if input.Address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Address is required"})
		return
	}

	// Generate EmployeeID format EMP-xxx
	var lastEmployee models.Employee
	if err := config.DB.Order("id desc").First(&lastEmployee).Error; err != nil && err.Error() != "record not found" {
		fmt.Println("DB error:", err.Error()) // log internal error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	nextID := 1
	if lastEmployee.ID > 0 {
		nextID = int(lastEmployee.ID) + 1
	}
	employeeID := fmt.Sprintf("EMP-%03d", nextID)

	employee := models.Employee{
		EmployeeID:   employeeID,
		DepartmentID: input.DepartmentID,
		Name:         input.Name,
		Address:      input.Address,
	}

	if err := config.DB.Create(&employee).Error; err != nil {
		fmt.Println("DB error:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Return with response struct
	if err := config.DB.Preload("Department").First(&employee, employee.ID).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{"data": toEmployeeDetailResp(employee)})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": employee})
	}
}

// Update employee by ID
func UpdateEmployee(c *gin.Context) {
	id := c.Param("id")
	var employee models.Employee
	if err := config.DB.First(&employee, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	var input struct {
		DepartmentID uint   `form:"department_id"`
		Name         string `form:"name"`
		Address      string `form:"address"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	// Custom validation
	if input.DepartmentID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Department is required"})
		return
	}
	if input.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}
	if input.Address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Address is required"})
		return
	}

	employee.DepartmentID = input.DepartmentID
	employee.Name = input.Name
	employee.Address = input.Address

	if err := config.DB.Save(&employee).Error; err != nil {
		fmt.Println("DB error:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Return with response struct
	if err := config.DB.Preload("Department").First(&employee, employee.ID).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{"data": toEmployeeDetailResp(employee)})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": employee})
	}
}

// Delete
func DeleteEmployee(c *gin.Context) {
	id := c.Param("id")
	var employee models.Employee
	if err := config.DB.First(&employee, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}
	// Delete attendance history & attendance
	config.DB.Where("employee_id = ?", employee.EmployeeID).Delete(&models.AttendanceHistory{})
	config.DB.Where("employee_id = ?", employee.EmployeeID).Delete(&models.Attendance{})
	// Delete employee
	if err := config.DB.Delete(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete employee"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Employee deleted successfully"})
}
