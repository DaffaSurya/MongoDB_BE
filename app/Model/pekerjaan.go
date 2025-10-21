package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Pekerjaan struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AlumniID        primitive.ObjectID `bson:"alumni_id,omitempty" json:"alumni_id,omitempty"`
	Nama_perusahaan string             `bson:"nama_perusahaan,omitempty" json:"nama_perusahaan,omitempty"`
	Posisi_jabatan  string             `bson:"posisi_jabatan,omitempty" json:"posisi_jabatan,omitempty"`
	Bidang_Industri string             `bson:"bidang_industri,omitempty" json:"bidang_industri,omitempty"`
	Lokasi_kerja    string             `bson: "lokasi_kerja,omitempty" json: "lokasi_kerja,omitempty"`
	Gaji_range      string             `bson: "gaji_range,omitempty" json: "gaji_range,omitempty"`
	Tanggal_Kerja   int64              `bson: "tanggal_kerja" json: "tanggal_kerja"`
	Tanggal_selesai int64              `bson: "tanggal_selesai" json: "tanggal_selesai"`
	Status          string             `bson: "status" json: "status"`
	Description     string             `bson:"deskripsi,omitempty" json:"deskripsi,omitempty"`
	CreatedAt       int64              `bson:"created_at" json:"created_at"`
	UpdatedAt       int64              `bson:"Updated_at" json:"Updated_at"`
}
