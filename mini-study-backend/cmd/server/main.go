// @title Mini Study API
// @version 0.1.0
// @description 企业学习平台后端 API 文档
// @termsOfService http://swagger.io/terms/
//
// @contact.name API Support
// @contact.email support@example.com
//
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host localhost:8080
// @BasePath /api/v1
// @schemes http https
//
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description 使用 "Bearer {token}" 格式，token 通过登录接口获取
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/javapub/mini-study/mini-study-backend/internal/bootstrap"
	"github.com/javapub/mini-study/mini-study-backend/internal/handler"
	"github.com/javapub/mini-study/mini-study-backend/internal/repository"
	"github.com/javapub/mini-study/mini-study-backend/internal/service"
)

func main() {
	cfg, err := bootstrap.LoadConfig("./configs")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	logger, err := bootstrap.InitLogger(cfg)
	if err != nil {
		log.Fatalf("init logger: %v", err)
	}
	defer logger.Sync() //nolint:errcheck

	db, err := bootstrap.InitDatabase(cfg, logger)
	if err != nil {
		logger.Fatal("init database", zap.Error(err))
	}

	validate := validator.New()

	auditRepo := repository.NewAuditRepository(db)
	userRepo := repository.NewUserRepository(db)
	relationRepo := repository.NewManagerEmployeeRepository(db)
	contentCategoryRepo := repository.NewContentCategoryRepository(db)
	contentRepo := repository.NewContentRepository(db)
	learningRecordRepo := repository.NewLearningRecordRepository(db)
	bannerRepo := repository.NewBannerRepository(db)

	auditService := service.NewAuditService(auditRepo)
	tokenService := service.NewTokenService(cfg.JWT.Secret, cfg.JWT.Issuer, cfg.JWT.TTL, cfg.JWT.RefreshTTL, userRepo)
	userService := service.NewUserService(userRepo, relationRepo, auditService)
	contentService := service.NewContentService(contentCategoryRepo, contentRepo, userRepo)
	learningService := service.NewLearningService(learningRecordRepo, contentRepo, userRepo)
	bannerService := service.NewBannerService(bannerRepo, userRepo, auditService)

	userHandler := handler.NewUserHandler(userService, tokenService)
	contentHandler := handler.NewContentHandler(contentService)
	learningHandler := handler.NewLearningHandler(learningService)
	bannerHandler := handler.NewBannerHandler(bannerService)
	uploadHandler := handler.NewUploadHandler(cfg.Upload.Dir, cfg.Upload.MaxSizeMB, auditService)
	systemHandler := handler.NewSystemHandler(cfg.App.Name, cfg.App.Version)

	engine := gin.New()
	bootstrap.RegisterMiddlewares(engine, cfg, logger, validate)
	bootstrap.RegisterRoutes(engine, cfg, userHandler, contentHandler, learningHandler, bannerHandler, uploadHandler, systemHandler)

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      engine,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		logger.Info("server starting", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server start failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.RequestTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server shutdown error", zap.Error(err))
	}

	logger.Info("server exited")
}
