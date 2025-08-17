package models

import (
	"time"
)

type Department struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	DepartmentName  string    `gorm:"type:varchar(255);not null" json:"department_name"`
	MaxClockInTime  string    `gorm:"type:time;not null" json:"max_clock_in_time"`
	MaxClockOutTime string    `gorm:"type:time;not null" json:"max_clock_out_time"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	Employees []Employee `gorm:"foreignKey:DepartmentID;references:ID"`
}

func (Department) TableName() string {
	return "departments"
}
