package repository

import (
	model "Mango/app/Model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	Col *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		Col: db.Collection("Users"),
	}
}

// ✅ Cari user berdasarkan ID (digunakan di AuthMiddleware)
func (r *UserRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	var user model.User
	err := r.Col.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	return &user, err
}

// ✅ Cari user berdasarkan username (untuk login)
func (r *UserRepository) FindByUsername(ctx context.Context, Username string) (*model.User, error) {
	var user model.User
	err := r.Col.FindOne(ctx, bson.M{"username": Username}).Decode(&user)
	return &user, err
}

// ✅ Tambah user baru (untuk register)
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	_, err := r.Col.InsertOne(ctx, user)
	return err
}
