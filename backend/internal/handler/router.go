package handler

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/Leander-Wendt/Penates/backend/internal/config"
	"github.com/Leander-Wendt/Penates/backend/internal/middleware"
	"github.com/Leander-Wendt/Penates/backend/internal/models"
	"github.com/Leander-Wendt/Penates/backend/internal/repository"
	"github.com/Leander-Wendt/Penates/backend/internal/service"
)

// NewRouter builds the complete Gin engine for the Penates API: it wires
// repositories, services, and handlers over db, and registers every route
// under /healthz, /api/v1, and /uploads.
func NewRouter(cfg *config.Config, db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS(cfg.CORSAllowedOrigins))

	orgRepo := repository.NewGormOrganisationRepository(db)
	userRepo := repository.NewGormUserRepository(db)
	itemRepo := repository.NewGormItemRepository(db)
	loanRequestRepo := repository.NewGormLoanRequestRepository(db)

	orgSvc := service.NewOrganisationService(orgRepo)
	userSvc := service.NewUserService(userRepo)
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpiry)
	itemSvc := service.NewItemService(itemRepo)
	loanRequestSvc := service.NewLoanRequestService(loanRequestRepo)
	uploadSvc := service.NewUploadService(itemRepo, cfg.UploadDir, cfg.MaxUploadSizeMB)

	orgHandler := NewOrganisationHandler(orgSvc)
	userHandler := NewUserHandler(userSvc)
	authHandler := NewAuthHandler(authSvc)
	itemHandler := NewItemHandler(itemSvc)
	loanRequestHandler := NewLoanRequestHandler(loanRequestSvc)
	uploadHandler := NewUploadHandler(uploadSvc)

	r.GET("/healthz", healthzHandler(db))

	auth := middleware.RequireAuth(cfg.JWTSecret)
	adminLogistics := middleware.RequireRoles(models.RoleAdmin, models.RoleLogistics)
	adminOnly := middleware.RequireRoles(models.RoleAdmin)

	r.GET("/uploads/:filename", auth, uploadsHandler(cfg.UploadDir))

	v1 := r.Group("/api/v1")
	{
		v1.POST("/auth/login", authHandler.Login)

		organisations := v1.Group("/organisations", auth)
		{
			organisations.GET("", orgHandler.List)
			organisations.POST("", adminLogistics, orgHandler.Create)
			organisations.GET("/:id", orgHandler.Get)
			organisations.PUT("/:id", adminLogistics, orgHandler.Update)
			organisations.DELETE("/:id", adminLogistics, orgHandler.Delete)
		}

		users := v1.Group("/users", auth)
		{
			users.GET("/me", userHandler.Me)
			users.GET("", adminLogistics, userHandler.List)
			users.POST("", adminLogistics, userHandler.Create)
			users.GET("/:id", userHandler.Get)
			users.PUT("/:id", adminLogistics, userHandler.Update)
			users.DELETE("/:id", adminOnly, userHandler.Delete)
		}

		items := v1.Group("/items", auth)
		{
			items.GET("", itemHandler.List)
			items.POST("", adminLogistics, itemHandler.Create)
			items.GET("/:inventoryNumber", itemHandler.Get)
			items.PUT("/:inventoryNumber", adminLogistics, itemHandler.Update)
			items.DELETE("/:inventoryNumber", adminLogistics, itemHandler.Delete)
			items.POST("/:inventoryNumber/image", adminLogistics, uploadHandler.Upload)
			items.DELETE("/:inventoryNumber/image", adminLogistics, uploadHandler.Delete)
		}

		loanRequests := v1.Group("/loan-requests", auth)
		{
			loanRequests.GET("", loanRequestHandler.List)
			loanRequests.POST("", loanRequestHandler.Create)
			loanRequests.GET("/:id", loanRequestHandler.Get)
			loanRequests.PUT("/:id", loanRequestHandler.Update)
			loanRequests.DELETE("/:id", loanRequestHandler.Delete)
			loanRequests.PATCH("/:id/status", adminLogistics, loanRequestHandler.UpdateStatus)
		}
	}

	return r
}

// healthzHandler returns a handler for GET /healthz that also verifies
// connectivity to the database.
func healthzHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

// uploadsHandler returns a handler for GET /uploads/:filename that serves
// files from uploadDir to any authenticated user.
func uploadsHandler(uploadDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		filename := c.Param("filename")
		path := filepath.Join(uploadDir, filepath.Base(filename))
		if _, err := os.Stat(path); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.File(path)
	}
}
