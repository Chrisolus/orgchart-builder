package models

type Role struct {
	ID   uint   `gorm:"primaryKey" json:"id" binding:"required"`
	Name string `gorm:"type:varchar(100);unique;not null" json:"name" binding:"required"`
}
