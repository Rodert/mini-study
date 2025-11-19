package main

import (
	"log"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/javapub/mini-study/mini-study-backend/internal/bootstrap"
	"github.com/javapub/mini-study/mini-study-backend/internal/model"
	"github.com/javapub/mini-study/mini-study-backend/internal/utils"
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
		logger.Fatal("migrate failed", zap.Error(err))
	}

	ensureDefaultAdmin(db, logger)
	ensureDefaultCategories(db, logger)
	ensureDefaultBanners(db, logger)

	logger.Info("database migrated")
}

// ensureDefaultAdmin creates a default admin user if none exists.
func ensureDefaultAdmin(db *gorm.DB, logger *zap.Logger) {
	const defaultWorkNo = "admin"
	const defaultPassword = "admin123456"

	var count int64
	if err := db.Model(&model.User{}).Where("role = ?", model.RoleAdmin).Count(&count).Error; err != nil {
		logger.Error("count admin users failed", zap.Error(err))
		return
	}
	if count > 0 {
		logger.Info("admin user already exists, skip seeding")
		return
	}

	hash, err := utils.HashPassword(defaultPassword)
	if err != nil {
		logger.Error("hash default admin password failed", zap.Error(err))
		return
	}

	admin := &model.User{
		WorkNo:       defaultWorkNo,
		Name:         "系统管理员",
		Phone:        "",
		PasswordHash: hash,
		Role:         model.RoleAdmin,
		Status:       true,
	}

	if err := db.Create(admin).Error; err != nil {
		logger.Error("create default admin user failed", zap.Error(err))
		return
	}

	logger.Info("default admin user created",
		zap.String("work_no", defaultWorkNo),
	)
}

// ensureDefaultCategories seeds default content categories.
func ensureDefaultCategories(db *gorm.DB, logger *zap.Logger) {
	defaults := []model.ContentCategory{
		{Name: "专题特训", RoleScope: model.RoleEmployee, SortOrder: 10, Status: true},
		{Name: "金牌课程", RoleScope: model.RoleEmployee, SortOrder: 20, Status: true},
		{Name: "新人养成", RoleScope: model.RoleEmployee, SortOrder: 30, Status: true},
		{Name: "短视频教学", RoleScope: model.RoleEmployee, SortOrder: 40, Status: true},
		{Name: "营销推广", RoleScope: model.RoleEmployee, SortOrder: 50, Status: true},
		{Name: "学习执行", RoleScope: model.RoleEmployee, SortOrder: 60, Status: true},
		{Name: "销售技巧学习", RoleScope: model.RoleEmployee, SortOrder: 70, Status: true},
		{Name: "新员工培训", RoleScope: model.RoleManager, SortOrder: 10, Status: true},
		{Name: "老员工进阶", RoleScope: model.RoleManager, SortOrder: 20, Status: true},
	}

	for _, item := range defaults {
		var count int64
		if err := db.Model(&model.ContentCategory{}).Where("name = ?", item.Name).Count(&count).Error; err != nil {
			logger.Error("count category failed", zap.Error(err), zap.String("name", item.Name))
			continue
		}
		if count > 0 {
			continue
		}
		if err := db.Create(&item).Error; err != nil {
			logger.Error("create category failed", zap.Error(err), zap.String("name", item.Name))
		} else {
			logger.Info("category seeded", zap.String("name", item.Name))
		}
	}
}

// ensureDefaultBanners seeds a sample banner for both roles.
func ensureDefaultBanners(db *gorm.DB, logger *zap.Logger) {
	const title = "欢迎来到学习中心"

	var count int64
	if err := db.Model(&model.Banner{}).Where("title = ?", title).Count(&count).Error; err != nil {
		logger.Error("count banner failed", zap.Error(err))
		return
	}
	if count > 0 {
		logger.Info("default banner already exists, skip")
		return
	}

	now := time.Now()
	end := now.AddDate(0, 3, 0)
	banner := &model.Banner{
		Title:        title,
		ImageURL:     "https://example.com/banners/default.png",
		LinkURL:      "https://example.com/h5/welcome",
		VisibleRoles: "both",
		SortOrder:    10,
		Status:       true,
		StartAt:      &now,
		EndAt:        &end,
	}

	if err := db.Create(banner).Error; err != nil {
		logger.Error("create default banner failed", zap.Error(err))
		return
	}
	logger.Info("default banner created", zap.String("title", title))
}
