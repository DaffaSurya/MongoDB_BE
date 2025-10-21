package routes

import (
	"Mango/app/service"
	"Mango/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.RouterGroup, authService *service.AuthService) {
	r.POST("/register", authService.Register)
	r.POST("/login", authService.Login)
}

func AlumniRoutes(r *gin.RouterGroup, alumniService *service.AlumniService) {
	alumni := r.Group("/alumni")
	{
		// 🔹 Semua user yang login bisa melihat daftar alumni
		alumni.GET("/", alumniService.GetAllAlumni)

		// 🔹 Hanya admin yang bisa menambah, update, dan hapus data alumni
		alumni.POST("/", middleware.RoleMiddleware("admin"), alumniService.CreateAlumni)
		alumni.PUT("/:id", middleware.RoleMiddleware("admin"), alumniService.UpdateAlumni)
		alumni.DELETE("/:id", middleware.RoleMiddleware("admin"), alumniService.DeleteAlumni)

		// 🔹 Detail alumni bisa dilihat siapa saja yang login
		alumni.GET("/:id", alumniService.GetAlumniByID)
	}
}

func PekerjaanRoutes(r *gin.RouterGroup, pekerjaanService *service.PekerjaanService) {
	pekerjaan := r.Group("/pekerjaan")
	{
		// 🔹 Semua user bisa melihat pekerjaan
		pekerjaan.GET("/", pekerjaanService.GetAllPekerjaan)
		pekerjaan.GET("/alumni_id", pekerjaanService.GetPekerjaanByAlumni)

		// 🔹 Hanya admin yang boleh create, update, dan delete pekerjaan
		pekerjaan.POST("/", middleware.RoleMiddleware("admin"), pekerjaanService.CreatePekerjaan)
		pekerjaan.PUT("/:id", middleware.RoleMiddleware("admin"), pekerjaanService.UpdatePekerjaan)
		pekerjaan.DELETE("/:id", middleware.RoleMiddleware("admin"), pekerjaanService.DeletePekerjaan)
	}
}
