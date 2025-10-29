package repository

import (
	model "Mango/app/Model"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Filerepository struct {
	Col *mongo.Collection
}


func NewUploadRepository(db *mongo.Database) *Filerepository {
	return &Filerepository{Col: db.Collection("uploads")}
}


func (r *Filerepository) Save(ctx context.Context, file *model.Files) error {
	_, err := r.Col.InsertOne(ctx, file)
	return err
}