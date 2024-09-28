package domain

import "time"

type Shop struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Email     string    `gorm:"uniqueIndex;type:varchar(255);unique;not null" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"`
	Name      string    `gorm:"uniqueIndex;type:varchar(255);not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
