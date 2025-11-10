// @title Mango API Documentation
// @version 1.0
// @description API untuk manajemen data alumni menggunakan MongoDB dan Gin
// @termsOfService http://swagger.io/terms/

// @contact.name Daffa Surya
// @contact.url https://github.com/daffasurya
// @contact.email daffa@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
package main

import (
	"Mango/app/repository"
	"Mango/app/service"
	"Mango/middleware"
	"Mango/routes"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "Mango/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// 🔹 Load .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("❌ Error loading .env file")
	}

	// 🔹 Ambil variabel dari .env
	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB")
	port := os.Getenv("SERVER_PORT")

	if mongoURI == "" || dbName == "" {
		log.Fatal("❌ MONGO_URI or MONGO_DB not found in .env")
	}

	// 🔹 Koneksi ke MongoDB
	clientOptions := options.Client().ApplyURI(mongoURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("❌ MongoDB connection error:", err)
	}

	// 🔹 Tes koneksi
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("❌ Cannot connect to MongoDB:", err)
	}
	fmt.Println("✅ Connected to MongoDB!")

	// 🔹 Inisialisasi database dan repository
	db := client.Database(dbName)
	userRepo := repository.NewUserRepository(db)
	alumniRepo := repository.NewAlumniRepository(db)
	pekerjaanRepo := repository.NewPekerjaanRepository(db)
	uploadRepo := repository.NewUploadRepository(db)

	// 🔹 Inisialisasi service
	authService := service.NewAuthService(userRepo)
	alumniService := service.NewAlumniService(alumniRepo)
	pekerjaanService := service.NewPekerjaanService(pekerjaanRepo)
	// FileService := service.FileService(FileRepo)
	uploadService := service.NewFileservice(uploadRepo)

	// 🔹 Setup router Gin
	router := gin.Default()
	router.Use(gin.Logger(), gin.Recovery())

	// =============================
	// 🔹 EndPoint Swagger
	// =============================

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	fmt.Println("📘 Swagger UI available at: http://localhost:" + port + "/swagger/index.html")

	// =============================
	// 🔹 ROUTING SECTION
	// =============================

	// 1️⃣ Public routes (tanpa middleware)
	public := router.Group("/auth")
	{
		routes.AuthRoutes(public, authService)
	}

	// 2️⃣ Protected routes (dengan middleware)
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware(userRepo))
	{
		routes.AlumniRoutes(api, alumniService)
		routes.PekerjaanRoutes(api, pekerjaanService)

	}

	// buat router untuk fitur uploads
	auth := router.Group("/uploads")
	auth.Use(middleware.AuthMiddleware(userRepo))
	{
		routes.FileRoutes(auth, uploadService, userRepo)
	}

	// 🔹 Jalankan server
	if port == "" {
		port = "3000"
	}
	fmt.Printf("🚀 Server running on port %s\n", port)
	router.Run(":" + port)

}
