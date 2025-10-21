package repository

import (
	model "Mango/app/Model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AlumniRepository struct {
	Col *mongo.Collection
}

func NewAlumniRepository(db *mongo.Database) *AlumniRepository {
	return &AlumniRepository{Col: db.Collection("alumni")}
}

func (r *AlumniRepository) GetAllAlumni(ctx context.Context) ([]model.Alumni, error) {
	// Buat slice untuk menampung semua data alumni
	var alumniList []model.Alumni

	// Lakukan query untuk mengambil semua dokumen di koleksi
	cursor, err := r.Col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Loop hasil query dan decode ke struct
	for cursor.Next(ctx) {
		var a model.Alumni
		if err := cursor.Decode(&a); err != nil {
			return nil, err
		}
		alumniList = append(alumniList, a)
	}

	// Cek apakah ada error selama iterasi cursor
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return alumniList, nil
}

func (r *AlumniRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Alumni, error) {
	var a model.Alumni
	err := r.Col.FindOne(ctx, bson.M{"_id": id}).Decode(&a)
	return &a, err
}

func (r *AlumniRepository) Create(ctx context.Context, alum *model.Alumni) error {
	alum.ID = primitive.NewObjectID()
	alum.CreatedAt = time.Now().Unix()
	_, err := r.Col.InsertOne(ctx, alum)
	return err
}

func (r *AlumniRepository) Update(ctx context.Context, id primitive.ObjectID, alum *model.Alumni) error {
	update := bson.M{
		"$set": bson.M{
			"nim":        alum.NIM,
			"nama":       alum.Nama,
			"jurusan":    alum.Jurusan,
			"updated_at": time.Now().Unix(),
		},
	}

	_, err := r.Col.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *AlumniRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	// Filter berdasarkan _id
	filter := bson.M{"_id": id}

	// Jalankan operasi DeleteOne
	result, err := r.Col.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	// Jika tidak ada dokumen yang dihapus
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
