package service

import (
	model "Mango/app/Model"
	"Mango/app/repository"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PekerjaanService struct {
	Repo *repository.PekerjaanRepository
}

func NewPekerjaanService(repo *repository.PekerjaanRepository) *PekerjaanService {
	return &PekerjaanService{Repo: repo}
}

// ✅ Create pekerjaan baru
func (s *PekerjaanService) CreatePekerjaan(c *gin.Context) {
	alumniID := c.MustGet("alumni_id").(primitive.ObjectID)

	var req struct {
		Nama            string `json:"nama_perusahaan" binding:"required"`
		Posisi          string `json:"posisi_jabatan" binding:"required"`
		Bidang_Industri string `json:"bidang_industri" binding:"required"`
		Lokasi          string `json:"lokasi_kerja" binding:"required"`
		TahunMasuk      int    `json:"tanggal_kerja" binding:"required"`
		TahunKeluar     int    `json:"tanggal_selesai"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pekerjaan := model.Pekerjaan{
		AlumniID:        alumniID,
		Nama_perusahaan: req.Nama,
		Posisi_jabatan:  req.Posisi,
		Bidang_Industri: req.Bidang_Industri,
		Lokasi_kerja:    req.Lokasi,
		Tanggal_Kerja:   int64(req.TahunMasuk),
		Tanggal_selesai: int64(req.TahunKeluar),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Repo.Create(ctx, &pekerjaan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create pekerjaan"})
		return
	}

	c.JSON(http.StatusCreated, pekerjaan)
}

func (s *PekerjaanService) GetAllPekerjaan(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	results, err := s.Repo.GetAllPekerjaan(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pekerjaan"})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (s *PekerjaanService) GetPekerjaanByAlumni(c *gin.Context) {
	// Ambil user dari context yang diset di middleware
	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	user := userVal.(*model.User)

	// Konversi user.ID ke alumniID jika diperlukan
	alumniID := user.ID

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	results, err := s.Repo.FindByAlumniID(ctx, alumniID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pekerjaan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data pekerjaan berhasil diambil",
		"data":    results,
	})
}


func (s *PekerjaanService) GetPekerjaanByID(c *gin.Context) {
	idParam := c.Param("id")

	// Convert ID dari string ke ObjectID MongoDB
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pekerjaan, err := s.Repo.FindByID(ctx, objID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Pekerjaan not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pekerjaan"})
		return
	}

	c.JSON(http.StatusOK, pekerjaan)
}


// ✅ Update pekerjaan tertentu
func (s *PekerjaanService) UpdatePekerjaan(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pekerjaan id"})
		return
	}

	var req struct {
		Nama            string `json:"nama_perusahaan"`
		Posisi          string `json:"posisi_jabatan"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{}
	if req.Nama != "" {
		update["nama_perusahaan"] = req.Nama
	}
	if req.Posisi != "" {
		update["posisi_jabatan"] = req.Posisi
	}
	

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Repo.Update(ctx, objID, update); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update pekerjaan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "pekerjaan updated"})
}

// ✅ Delete pekerjaan
func (s *PekerjaanService) DeletePekerjaan(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pekerjaan id"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Repo.Delete(ctx, objID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete pekerjaan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "pekerjaan deleted"})
}
