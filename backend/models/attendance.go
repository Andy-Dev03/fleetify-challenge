package models

import "time"

type Attendance struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	EmployeeID   string     `gorm:"type:varchar(50);not null" json:"employee_id"`
	AttendanceID string     `gorm:"type:varchar(100);not null" json:"attendance_id"`
	ClockIn      time.Time  `gorm:"type:timestamp" json:"clock_in"`
	ClockOut     *time.Time `gorm:"type:timestamp" json:"clock_out"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	Employee          Employee `gorm:"foreignKey:EmployeeID;references:EmployeeID"`
}
