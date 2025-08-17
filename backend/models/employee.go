package models

import (
	"time"
)

type Employee struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	EmployeeID   string    `gorm:"unique;type:varchar(50)" json:"employee_id"`
	DepartmentID uint      `gorm:"column:department_id;not null" json:"department_id"`
	Name         string    `gorm:"type:varchar(255)" json:"name"`
	Address      string    `gorm:"type:text" json:"address"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Department        Department `gorm:"foreignKey:DepartmentID;references:ID"`
}
