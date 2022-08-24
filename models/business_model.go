package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Business struct {
	Id           primitive.ObjectID `json:"id,omitempty"`
	BusinessName string             `json:"biz_name,omitempty" validate:"required"`
	BusinessId   string             `json:"biz_id,omitempty" validate:"required"`
	CreatedAt    time.Time          `json:"created_at,omitempty"`
	UpdatedAt    time.Time          `json:"updated_at,omitempty"`
}
