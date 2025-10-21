package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string             `bson:"username" json:"username"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"password"`
	Role      string             `bson:"role" json:"role"`
	AlumniID  primitive.ObjectID `bson:"alumni_id,omitempty" json:"alumni_id,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
