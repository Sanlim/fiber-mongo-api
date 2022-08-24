package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Business struct {
	Id            primitive.ObjectID `json:"id,omitempty"`
	Business_Name string             `json:"biz_name,omitempty" validate:"required"`
	Business_Id   string             `json:"biz_id,omitempty" validate:"required"`
	Created_At    time.Time          `json:"created_at,omitempty"`
	Updated_At    time.Time          `json:"updated_at,omitempty"`
}
