package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Files struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID primitive.ObjectID `bson: "user_id,omitempty" json:"user_id"`
	Type string `bson:"type" json:"type"`
	Filepath string `bson:"file_path" json:"file_path"`
	Filename string `bson:"file_name"  json:"file_name"`
	ContentType string   `bson:"content_type" json:"content_type"`
	UploadedAt time.Time `bson:"uploaded_at" json:"uploaded_at"`
}


