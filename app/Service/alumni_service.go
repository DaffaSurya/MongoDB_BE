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

type AlumniService struct {
	repo *repository.AlumniRepository
}

func NewAlumniService(r *repository.AlumniRepository) *AlumniService {
	return &AlumniService{repo: r}
}

// 🟢 Get all alumni
func (s *AlumniService) GetAllAlumni(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	alumni, err := s.repo.GetAllAlumni(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, alumni)
}

// 🟢 Get alumni by ID
func (s *AlumniService) GetAlumniByID(c *gin.Context) {
	id := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	alumni, err := s.repo.FindByID(ctx, objID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Alumni not found"})
		return
	}

	c.JSON(http.StatusOK, alumni)
}

// 🟢 Create new alumni
func (s *AlumniService) CreateAlumni(c *gin.Context) {
	var alum model.Alumni
	if err := c.ShouldBindJSON(&alum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.repo.Create(ctx, &alum); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create alumni"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Data alumni berhasil ditambahkan",
		"data":    alum,
	})
}

func (s *AlumniService) UpdateAlumni(c *gin.Context) {
	idParam := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var updatedData model.Alumni
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.repo.Update(ctx, objID, &updatedData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update alumni"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data alumni berhasil diperbarui",
		"data":    updatedData,
	})
}

// 🟢 Delete alumni
func (s *AlumniService) DeleteAlumni(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Delete pakai DeleteOne (karena repo kamu belum punya delete, bisa ditambah)
	_, err = s.repo.Col.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete alumni"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Alumni deleted successfully"})
}
