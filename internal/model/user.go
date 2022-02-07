package model

import uuid "github.com/gofrs/uuid"

type User struct {
	ID       uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id; default:uuid_generate_v4()" json:"id,omitempty" `
	Username string    `gorm:"type:varchar(150); not null; unique" json:"name" binding:"required"`
	Email    string    `gorm:"type:varchar(150); not null; unique" json:"email" binding:"required" valid:"email"`
	Phone    string    `gorm:"type:varchar(150); not null; unique" json:"phone" binding:"required" valid:"numeric"`
	Password string    `gorm:"type:varchar(150); not null; unique" json:"password" binding:"required"`
}
