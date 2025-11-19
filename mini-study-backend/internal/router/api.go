package router

import (
    "github.com/gin-gonic/gin"

    "github.com/javapub/mini-study/mini-study-backend/internal/handler"
)

// RegisterRoutes wires API and system routes.
func RegisterRoutes(
	engine *gin.Engine,
	swaggerEnabled bool,
	authMiddleware gin.HandlerFunc,
	userHandler *handler.UserHandler,
	contentHandler *handler.ContentHandler,
	learningHandler *handler.LearningHandler,
	bannerHandler *handler.BannerHandler,
	uploadHandler *handler.UploadHandler,
	systemHandler *handler.SystemHandler,
) {
    api := engine.Group("/api/v1")

    user := api.Group("/users")
    {
        user.POST("/register", userHandler.Register)
        user.POST("/login", userHandler.Login)
		user.POST("/token/refresh", userHandler.RefreshToken)
		user.GET("/managers", userHandler.ListManagers)
	}

	// Authed user endpoints
	authUser := api.Group("/users")
	authUser.Use(authMiddleware)
	{
		authUser.PATCH("/me/profile", userHandler.UpdateProfile)
	}

	// Content routes (need auth)
	content := api.Group("/contents")
	content.Use(authMiddleware)
	{
		content.GET("/categories", contentHandler.ListCategories)
		content.GET("/", contentHandler.ListPublishedContents)
		content.GET("/:id", contentHandler.GetContentDetail)
	}

	// Learning routes
	learning := api.Group("/learning")
	learning.Use(authMiddleware)
	{
		learning.GET("/", learningHandler.ListProgress)
		learning.GET("/:content_id", learningHandler.GetProgress)
		learning.POST("/", learningHandler.UpdateProgress)
	}

	// Banner routes
	banners := api.Group("/banners")
	banners.Use(authMiddleware)
	{
		banners.GET("/", bannerHandler.ListVisibleBanners)
	}

	// Admin-only endpoints (admin check in handler/service)
	admin := api.Group("/admin")
	admin.Use(authMiddleware)
	{
		admin.POST("/managers", userHandler.AdminCreateManager)
		admin.POST("/users/:id/promote-manager", userHandler.AdminPromoteToManager)
		admin.PUT("/users/:id/managers", userHandler.AdminUpdateEmployeeManagers)

		adminContents := admin.Group("/contents")
		{
			adminContents.GET("/", contentHandler.AdminListContents)
			adminContents.POST("/", contentHandler.AdminCreateContent)
			adminContents.PUT("/:id", contentHandler.AdminUpdateContent)
		}

		adminBanners := admin.Group("/banners")
		{
			adminBanners.GET("/", bannerHandler.AdminListBanners)
			adminBanners.POST("/", bannerHandler.AdminCreateBanner)
			adminBanners.PUT("/:id", bannerHandler.AdminUpdateBanner)
		}
    }

    files := api.Group("/files")
	files.Use(authMiddleware)
    files.POST("/upload", uploadHandler.Upload)

    RegisterSystemRoutes(engine, systemHandler)

    if swaggerEnabled {
        RegisterSwagger(engine)
    }
}
