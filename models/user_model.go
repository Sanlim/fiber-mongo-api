package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id         primitive.ObjectID `json:"id,omitempty"`
	Name       string             `json:"name,omitempty" validate:"required"`
	Location   string             `json:"location,omitempty" validate:"required"`
	Title      string             `json:"title,omitempty" validate:"required"`
	Created_At time.Time          `json:"created_at,omitempty"`
	Updated_At time.Time          `json:"updated_at,omitempty"`
}

type UserPass struct {
	Id            primitive.ObjectID `json:"id,omitempty"`
	Username      string             `json:"username,omitempty" validate:"required"`
	Password      string             `json:"password,omitempty" validate:"required"`
	Business_Name string             `json:"biz_name,omitempty"`
	Business_Id   string             `json:"biz_id,omitempty"`
	Created_At    time.Time          `json:"created_at,omitempty"`
	Updated_At    time.Time          `json:"updated_at,omitempty"`
}
