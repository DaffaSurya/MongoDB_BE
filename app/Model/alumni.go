package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Alumni struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	NIM         string             `bson:"nim,omitempty" json:"nim,omitempty"`
	Nama        string             `bson:"nama" json:"nama"`
	Jurusan     string             `bson:"jurusan,omitempty" json:"jurusan,omitempty"`
	Angkatan    int                `bson:"angkatan,omitempty" json:"angkatan,omitempty"`
	Tahun_lulus int                `bson:"tahun_lulus,omitempty" json:"tahun_lulus,omitempty"`
	Email       string             `bson:"email" json:"email"`
	No_telp     string             `bson:"no_telp" json:"no_telp"`
	Alamat      string             `bson:"alamat" json:"alamat"`
	CreatedAt   int64              `bson:"created_at" json:"created_at"`
	UpdateddAt  int64              `bson:"Update_at" json:"Update_at"`
}

type AlumniResponse struct {
	Data []Alumni `json: "data"`
	Meta MetaInfo `json: "meta"`
}