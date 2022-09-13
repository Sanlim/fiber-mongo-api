package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	Id                  primitive.ObjectID `json:"id,omitempty"`
	Product_Code        string             `json:"product_code,omitempty" validate:"required"`
	Product_Name        string             `json:"product_name,omitempty" validate:"required"`
	Product_Type        string             `json:"product_type,omitempty" validate:"required"`
	Product_Unit        string             `json:"product_unit,omitempty" validate:"required"`
	Product_Amount      int64              `json:"product_amount,omitempty" validate:"required"`
	Product_Price       int64              `json:"product_price,omitempty" validate:"required"`
	Product_Description string             `json:"product_description,omitempty" validate:"required"`
	Created_At          time.Time          `json:"created_at,omitempty"`
	Updated_At          time.Time          `json:"updated_at,omitempty"`
}

type HistoryPurchased struct {
	Id                 primitive.ObjectID `json:"id,omitempty"`
	Product_Code       string             `json:"product_code,omitempty" validate:"required"`
	Product_Name       string             `json:"product_name,omitempty" validate:"required"`
	Product_Price      int64              `json:"product_price,omitempty" validate:"required"`
	Customer_Purchased []Customer         `json:"customer_purchased,omitempty" validate:"required"`
	Created_At         time.Time          `json:"created_at,omitempty"`
	Updated_At         time.Time          `json:"updated_at,omitempty"`
}

type Customer struct {
	Id               primitive.ObjectID `json:"id,omitempty"`
	Customer_Id      string             `json:"customer_id,omitempty" validate:"required"`
	Customer_Name    string             `json:"customer_name,omitempty" validate:"required"`
	Purchased_Date   int64              `json:"purchased_date,omitempty" validate:"required"`
	Purchased_Amount int64              `json:"purchased_amount,omitempty" validate:"required"`
	Purchased_Total  int64              `json:"purchased_total,omitempty" validate:"required"`
	Created_At       time.Time          `json:"created_at,omitempty"`
	Updated_At       time.Time          `json:"updated_at,omitempty"`
}
