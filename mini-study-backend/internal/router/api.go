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
	examHandler *handler.ExamHandler,
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
		authUser.GET("/me", userHandler.GetCurrentUser)
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
		learning.GET("/stats", learningHandler.GetUserStats)
		learning.GET("/:content_id", learningHandler.GetProgress)
		learning.GET("/content/:content_id/stats", learningHandler.GetContentStats)
		learning.POST("/", learningHandler.UpdateProgress)
	}

	// Exam routes
	exams := api.Group("/exams")
	exams.Use(authMiddleware)
	{
		exams.GET("/my/results", examHandler.ListMyResults)
		exams.GET("/", examHandler.ListAvailable)
		exams.GET("/:id", examHandler.GetExamDetail)
		exams.POST("/:id/submit", examHandler.SubmitExam)
	}

	manager := api.Group("/manager")
	manager.Use(authMiddleware)
	{
		manager.GET("/exams/overview", examHandler.ManagerOverview)
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
		admin.GET("/users", userHandler.AdminListUsers)
		admin.GET("/users/:id", userHandler.AdminGetUser)
		admin.PUT("/users/:id/role", userHandler.AdminUpdateUserRole)
		admin.POST("/managers", userHandler.AdminCreateManager)
		admin.POST("/employees", userHandler.AdminCreateEmployee)
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

		adminExams := admin.Group("/exams")
		{
			adminExams.GET("/", examHandler.AdminListExams)
			adminExams.GET("/:id", examHandler.AdminGetExam)
			adminExams.POST("/", examHandler.AdminCreateExam)
			adminExams.PUT("/:id", examHandler.AdminUpdateExam)
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
