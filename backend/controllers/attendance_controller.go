package controllers

import (
	"fleetify-backend/config"
	"fleetify-backend/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AttendanceResp struct {
	ID           uint       `json:"id"`
	EmployeeID   string     `json:"employee_id"`
	AttendanceID string     `json:"attendance_id"`
	ClockIn      time.Time  `json:"clock_in"`
	ClockOut     *time.Time `json:"clock_out"`
}

type AttendanceLogResp struct {
	ID             uint   `json:"id"`
	EmployeeID     string `json:"employee_id"`
	AttendanceID   string `json:"attendance_id"`
	Name           string `json:"name"`
	DateAttendance string `json:"date_attendance"`
	AttendanceType int    `json:"attendance_type"`
	Description    string `json:"description"`
	Department     string `json:"department"`
	ClockIn        string `json:"clock_in"`
	ClockOut       string `json:"clock_out"`
}

func GetAttendanceLogs(c *gin.Context) {
	dateParam := c.Query("date")
	departmentParam := c.Query("department_id")

	var histories []models.AttendanceHistory
	db := config.DB.
		Preload("Employee").
		Preload("Employee.Department").
		Preload("Attendance")

	// Filter tanggal (YYYY-MM-DD)
	if dateParam != "" {
		if t, err := time.Parse("2006-01-02", dateParam); err == nil {
			start := t
			end := t.Add(24 * time.Hour)
			db = db.Where("date_attendance >= ? AND date_attendance < ?", start, end)
		}
	}

	// Filter department
	if departmentParam != "" {
		db = db.Joins("JOIN employees ON attendance_histories.employee_id = employees.employee_id").
			Where("employees.department_id = ?", departmentParam)
	}

	// Ambil data
	if err := db.Find(&histories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Bentuk response
	var logs []AttendanceLogResp
	for _, history := range histories {
		attendance := history.Attendance

		clockIn := ""
		clockOut := ""
		if !attendance.ClockIn.IsZero() {
			clockIn = attendance.ClockIn.Format("15:04:05")
		}
		if attendance.ClockOut != nil {
			clockOut = attendance.ClockOut.Format("15:04:05")
		}

		empName := history.Employee.Name
		if empName == "" {
			empName = history.EmployeeID
		}
		deptName := history.Employee.Department.DepartmentName
		if deptName == "" {
			deptName = "-"
		}

		// Tentukan status absensi
		description := history.Description
		if history.AttendanceType == 1 { // Clock In
			maxIn := history.Employee.Department.MaxClockInTime
			if clockIn != "" && clockIn <= maxIn {
				description = "On Time (Check-in)"
			} else if clockIn != "" {
				description = "Late (Check-in)"
			}
		} else if history.AttendanceType == 2 { // Clock Out
			maxOut := history.Employee.Department.MaxClockOutTime
			if clockOut != "" && clockOut >= maxOut {
				description = "On Time (Check-out)"
			} else if clockOut != "" {
				description = "Early Leave"
			}
		}

		logs = append(logs, AttendanceLogResp{
			ID:             history.ID,
			EmployeeID:     history.EmployeeID,
			AttendanceID:   history.AttendanceID,
			Name:           empName,
			DateAttendance: history.DateAttendance.Format("2006-01-02 15:04:05"),
			AttendanceType: history.AttendanceType,
			Description:    description,
			Department:     deptName,
			ClockIn:        clockIn,
			ClockOut:       clockOut,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": logs})
}

// CreateAttendance membuat data absensi baru (clock in)
func CreateAttendance(c *gin.Context) {
	var input struct {
		EmployeeID string `form:"employee_id" json:"employee_id"`
		ClockIn    string `form:"clock_in" json:"clock_in"` // format: 2006-01-02 15:04:05
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	if input.EmployeeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Employee is required"})
		return
	}
	if input.ClockIn == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Clock In is required"})
		return
	}

	clockInTime, err := time.Parse("2006-01-02 15:04:05", input.ClockIn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid format for clock_in, expected YYYY-MM-DD HH:mm:ss"})
		return
	}

	// Generate AttendanceID otomatis
	var lastAttendance models.Attendance
	if err := config.DB.Order("id desc").First(&lastAttendance).Error; err != nil && err.Error() != "record not found" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get last attendance"})
		return
	}
	nextID := 1
	if lastAttendance.ID > 0 {
		nextID = int(lastAttendance.ID) + 1
	}
	attendanceID := fmt.Sprintf("ATT-%03d", nextID)

	attendance := models.Attendance{
		EmployeeID:   input.EmployeeID,
		AttendanceID: attendanceID,
		ClockIn:      clockInTime,
		ClockOut:     nil,
	}
	if err := config.DB.Create(&attendance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Simpan riwayat absensi
	history := models.AttendanceHistory{
		EmployeeID:     input.EmployeeID,
		AttendanceID:   attendanceID,
		DateAttendance: clockInTime,
		AttendanceType: 1,
		Description:    "On Time (Check-in)",
	}
	config.DB.Create(&history)

	resp := AttendanceResp{
		ID:           attendance.ID,
		EmployeeID:   attendance.EmployeeID,
		AttendanceID: attendance.AttendanceID,
		ClockIn:      attendance.ClockIn,
		ClockOut:     attendance.ClockOut,
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// UpdateAttendance
func UpdateAttendance(c *gin.Context) {
	attendanceID := c.Param("id")

	var input struct {
		ClockOut string `form:"clock_out" json:"clock_out"`
	}

	// Bind input
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	// Custom validation
	if input.ClockOut == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Clock Out is required"})
		return
	}

	// Cari attendance berdasarkan attendance_id
	var attendance models.Attendance
	if err := config.DB.Where("attendance_id = ?", attendanceID).First(&attendance).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attendance not found"})
		return
	}

	// Validasi format datetime
	clockOutTime, err := time.Parse("2006-01-02 15:04:05", input.ClockOut)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid format for clock_out, expected YYYY-MM-DD HH:mm:ss"})
		return
	}
	attendance.ClockOut = &clockOutTime

	// Update DB
	if err := config.DB.Save(&attendance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Simpan riwayat clock out
	history := models.AttendanceHistory{
		EmployeeID:     attendance.EmployeeID,
		AttendanceID:   attendance.AttendanceID,
		DateAttendance: clockOutTime,
		AttendanceType: 2,
		Description:    "On Time (Check-out)",
	}
	config.DB.Create(&history)

	// Response sederhana
	resp := AttendanceResp{
		ID:           attendance.ID,
		EmployeeID:   attendance.EmployeeID,
		AttendanceID: attendance.AttendanceID,
		ClockIn:      attendance.ClockIn,
		ClockOut:     attendance.ClockOut,
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}
