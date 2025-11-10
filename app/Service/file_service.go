package service

import (
	model "Mango/app/Model"
	"Mango/app/repository"
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileService struct {
	Repo *repository.Filerepository
}

func NewFileservice(repo *repository.Filerepository) *FileService {
	return &FileService{Repo: repo}
}

// ========================= UPLOAD FOTO =========================

// @Summary Upload user photo
// @Description Upload a profile photo (JPG, JPEG, PNG only, max 1MB)
// @Tags Uploads
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Photo file"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /uploads/photo [post]
func (s *FileService) UploadPhoto(c *gin.Context) {
	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"error": "User not found in context"})
		return
	}

	var role string
	var targetUserID string

	// Cek apakah user bertipe *model.User
	if u, ok := userVal.(*model.User); ok {
		role = u.Role
		targetUserID = fmt.Sprintf("%v", u.ID) // konversi ke string bila perlu
	} else if uMap, ok := userVal.(map[string]interface{}); ok {
		role = fmt.Sprintf("%v", uMap["role"])
		targetUserID = fmt.Sprintf("%v", uMap["id"])
	} else {
		c.JSON(500, gin.H{"error": "Invalid user type in context"})
		return
	}

	// Jika role admin, izinkan override user_id lewat query param
	if role == "admin" {
		if q := c.Query("user_id"); q != "" {
			targetUserID = q
		}
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file not found"})
		return
	}

	// Validasi format
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only jpg, jpeg, png allowed"})
		return
	}

	// Validasi ukuran file (max 1MB)
	if file.Size > 1*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large (max 1MB)"})
		return
	}

	saveDir := "uploads/photos"
	os.MkdirAll(saveDir, os.ModePerm)
	filePath := fmt.Sprintf("%s/%s_%s", saveDir, targetUserID, file.Filename)
	c.SaveUploadedFile(file, filePath)

	objID, _ := primitive.ObjectIDFromHex(targetUserID)
	upload := model.Files{
		ID:          primitive.NewObjectID(),
		UserID:      objID,
		Type:        "image",
		Filename:    file.Filename,
		Filepath:    filePath,
		ContentType: file.Header.Get("Content-Type"),
		UploadedAt:  time.Now(),
	}

	if err := s.Repo.Save(context.Background(), &upload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save photo metadata"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "photo uploaded successfully", "data": upload})
}

// ========================= UPLOAD SERTIFIKAT =========================

// @Summary Upload user certificate
// @Description Upload a certificate file (PDF only, max 2MB)
// @Tags Uploads
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Certificate file"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /uploads/certificate [post]
func (s *FileService) UploadCertificate(c *gin.Context) {

	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"error": "User not found in context"})
		return
	}

	var role string
	var targetUserID string

	// Cek apakah user bertipe *model.User
	if u, ok := userVal.(*model.User); ok {
		role = u.Role
		targetUserID = fmt.Sprintf("%v", u.ID) // konversi ke string bila perlu
	} else if uMap, ok := userVal.(map[string]interface{}); ok {
		role = fmt.Sprintf("%v", uMap["role"])
		targetUserID = fmt.Sprintf("%v", uMap["id"])
	} else {
		c.JSON(500, gin.H{"error": "Invalid user type in context"})
		return
	}

	// Jika role admin, izinkan override user_id lewat query param
	if role == "admin" {
		if q := c.Query("user_id"); q != "" {
			targetUserID = q
		}
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file not found"})
		return
	}

	ext := filepath.Ext(file.Filename)
	if ext != ".pdf" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only pdf allowed"})
		return
	}

	if file.Size > 2*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large (max 2MB)"})
		return
	}

	saveDir := "uploads/certificates"
	os.MkdirAll(saveDir, os.ModePerm)
	filePath := fmt.Sprintf("%s/%s_%s", saveDir, targetUserID, file.Filename)
	c.SaveUploadedFile(file, filePath)

	objID, _ := primitive.ObjectIDFromHex(targetUserID)
	upload := model.Files{ 
		ID:          primitive.NewObjectID(),
		UserID:      objID,
		Type:        "certificate",
		Filename:    file.Filename,
		Filepath:    filePath,
		ContentType: file.Header.Get("Content-Type"),
		UploadedAt:  time.Now(),
	}

	if err := s.Repo.Save(context.Background(), &upload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save certificate metadata"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "certificate uploaded successfully", "data": upload})
}
