package repository

import (
	model "Mango/app/Model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PekerjaanRepository struct {
	Col *mongo.Collection
}

func NewPekerjaanRepository(db *mongo.Database) *PekerjaanRepository {
	return &PekerjaanRepository{Col: db.Collection("pekerjaan_alumni")}
}

func (r *PekerjaanRepository) GetAllPekerjaan(ctx context.Context) ([]model.Pekerjaan, error) {
	var results []model.Pekerjaan
	cursor, err := r.Col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var p model.Pekerjaan
		if err := cursor.Decode(&p); err != nil {
			return nil, err
		}
		results = append(results, p)
	}
	return results, cursor.Err()
}


// Create pekerjaan
func (r *PekerjaanRepository) Create(ctx context.Context, p *model.Pekerjaan) error {
	p.ID = primitive.NewObjectID()
	_, err := r.Col.InsertOne(ctx, p)
	return err
}

func (r *PekerjaanRepository) FindByAlumniID(ctx context.Context, alumniID primitive.ObjectID) ([]model.Pekerjaan, error) {
	var results []model.Pekerjaan
	cursor, err := r.Col.Find(ctx, bson.M{"alumni_id": alumniID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var p model.Pekerjaan
		if err := cursor.Decode(&p); err != nil {
			return nil, err
		}
		results = append(results, p)
	}
	return results, cursor.Err()
}


// Get Pekerjaan By ID
func (r *PekerjaanRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Pekerjaan, error) {
	var pekerjaan model.Pekerjaan

	err := r.Col.FindOne(ctx, bson.M{"_id": id}).Decode(&pekerjaan)
	if err != nil {
		return nil, err
	}

	return &pekerjaan, nil
}


// Update pekerjaan
func (r *PekerjaanRepository) Update(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	_, err := r.Col.UpdateByID(ctx, id, bson.M{"$set": update})
	return err
}

// Delete pekerjaan
func (r *PekerjaanRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.Col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
