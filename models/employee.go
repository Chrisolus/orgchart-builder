package models

import (
	"errors"
	"org_chart/initializers"

	"gorm.io/gorm"
)

type Employee struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName string `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName  string `gorm:"type:varchar(100);not null" json:"last_name"`
	RoleID    uint   `json:"role_id"`
	Role      *Role  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ManagerID *uint  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"manager_id"`
	Manager   *Employee
}

type EmployeeSummary struct {
	ID   uint   `json:"id"`
	Name string `json:"name"` // Firstname + Lastname
}

type RequestBody struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	RoleID    uint   `json:"role_id" binding:"required"`
	ManagerID *uint  `json:"manager_id"`
}

func (r *RequestBody) Validate() error {
	if r.FirstName == "" || r.LastName == "" {
		return errors.New("name cannot be empty")
	}

	var count int64
	if initializers.DB.Model(&Role{}).Where("id = ?", r.RoleID).Count(&count); count == 0 {
		return errors.New("role doesn't exist")
	}
	if initializers.DB.Model(&Employee{}).Where("id = ?", r.ManagerID).Count(&count); r.ManagerID != nil && count == 0 {
		r.ManagerID = nil
	}

	return nil
}

// func (r *RequestBody) NormalizeManagerID() {
// 	if r.ManagerID != nil && *r.ManagerID == 0 {
// 		r.ManagerID = nil
// 	}
// }
