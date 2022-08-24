package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id        primitive.ObjectID `json:"id,omitempty"`
	Name      string             `json:"name,omitempty" validate:"required"`
	Location  string             `json:"location,omitempty" validate:"required"`
	Title     string             `json:"title,omitempty" validate:"required"`
	CreatedAt time.Time          `json:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty"`
}

type UserPass struct {
	Id           primitive.ObjectID `json:"id,omitempty"`
	Username     string             `json:"username,omitempty" validate:"required"`
	Password     string             `json:"password,omitempty" validate:"required"`
	BusinessName string             `json:"biz_name,omitempty"`
	BusinessId   string             `json:"biz_id,omitempty"`
	CreatedAt    time.Time          `json:"created_at,omitempty"`
	UpdatedAt    time.Time          `json:"updated_at,omitempty"`
}
