package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)
type User struct{
	ID							primitive.ObjectID       `bson:"_id"`
	First_name
	Last_name
	Password
	Email
	Phone
	Token
	User_type
	Refersh_token
	Created_at
	Updated_at
	User_id
}