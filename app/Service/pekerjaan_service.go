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


// ✅ Get semua pekerjaan milik alumni yang login
func (s *PekerjaanService) GetPekerjaanByAlumni(c *gin.Context) {
	alumniID := c.MustGet("alumni_id").(primitive.ObjectID)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	results, err := s.Repo.FindByAlumniID(ctx, alumniID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch pekerjaan"})
		return
	}

	c.JSON(http.StatusOK, results)
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
