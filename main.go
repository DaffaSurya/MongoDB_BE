// package main

// import (
// 	middleware "Mango/Middleware"
// 	"Mango/app/repository"
// 	"Mango/app/service"
// 	"Mango/routes"
// 	"context"
// 	"fmt"
// 	"log"
// 	"os"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/joho/godotenv"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// func main() {
// 	// 🔹 Load .env file
// 	if err := godotenv.Load(); err != nil {
// 		log.Fatal("❌ Error loading .env file")
// 	}

// 	// 🔹 Ambil variabel environment
// 	mongoURI := os.Getenv("MONGO_URI")
// 	dbName := os.Getenv("MONGO_DB")
// 	port := os.Getenv("SERVER_PORT")

// 	if mongoURI == "" || dbName == "" {
// 		log.Fatal("❌ MONGO_URI or MONGO_DB not found in .env")
// 	}

// 	// 🔹 Koneksi ke MongoDB
// 	clientOptions := options.Client().ApplyURI(mongoURI)
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	client, err := mongo.Connect(ctx, clientOptions)
// 	if err != nil {
// 		log.Fatal("❌ MongoDB connection error:", err)
// 	}

// 	// 🔹 Tes koneksi Mongo
// 	if err := client.Ping(ctx, nil); err != nil {
// 		log.Fatal("❌ Cannot connect to MongoDB:", err)
// 	}
// 	fmt.Println("✅ Connected to MongoDB!")

// 	// 🔹 Inisialisasi Database
// 	db := client.Database(dbName)

// 	// 🔹 Inisialisasi Repository
// 	userRepo := repository.NewUserRepository(db)
// 	alumniRepo := repository.NewAlumniRepository(db)
// 	pekerjaanRepo := repository.NewPekerjaanRepository(db)

// 	// 🔹 Inisialisasi Service
// 	authService := service.NewAuthService(userRepo)
// 	alumniService := service.NewAlumniService(alumniRepo)
// 	pekerjaanService := service.NewPekerjaanService(pekerjaanRepo)

// 	// 🔹 Setup Router Gin
// 	router := gin.Default()

// 	// 🔹 Gunakan Middleware
// 	router.Use(gin.Logger())
// 	router.Use(gin.Recovery())

// 	// 🔹 Auth Routes (tanpa middleware)
// 	// .AuthRoutes(router, authService)

// 	// 🔹 Alumni Routes (dengan middleware)
// 	auth := router.Group("/api")
// 	auth.Use(middleware.AuthMiddleware(userRepo))
// 	{
// 		routes.AlumniRoutes(auth, authService)
// 		routes.PekerjaanRoutes(auth, pekerjaanService)
// 	}

// 	// 🔹 Jalankan server
// 	if port == "" {
// 		port = "3000"
// 	}
// 	fmt.Printf("🚀 Server running on port %s\n", port)
// 	router.Run(":" + port)
// }


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

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	// 🔹 Inisialisasi service
	authService := service.NewAuthService(userRepo)
	alumniService := service.NewAlumniService(alumniRepo)
	pekerjaanService := service.NewPekerjaanService(pekerjaanRepo)

	// 🔹 Setup router Gin
	router := gin.Default()
	router.Use(gin.Logger(), gin.Recovery())

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

	// 🔹 Jalankan server
	if port == "" {
		port = "3000"
	}
	fmt.Printf("🚀 Server running on port %s\n", port)
	router.Run(":" + port)
}
