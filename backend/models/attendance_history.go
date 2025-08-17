package models

import (
	"time"
)

type AttendanceHistory struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	EmployeeID     string    `gorm:"type:varchar(50);not null" json:"employee_id"`
	AttendanceID   string    `gorm:"type:varchar(100);not null" json:"attendance_id"`
	DateAttendance time.Time `gorm:"type:timestamp" json:"date_attendance"`
	AttendanceType int       `gorm:"type:tinyint" json:"attendance_type"` // 1=In, 2=Out
	Description    string    `gorm:"type:text;" json:"description"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	Employee   Employee   `gorm:"foreignKey:EmployeeID;references:EmployeeID"`
	Attendance Attendance `gorm:"foreignKey:AttendanceID;references:AttendanceID"`
}
