package bootstrap

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/javapub/mini-study/mini-study-backend/internal/model"
)

// InitDatabase configures the GORM connection and performs automigration.
func InitDatabase(cfg *Config, logger *zap.Logger) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch cfg.Database.Driver {
	case "sqlite", "sqlite3":
		dialector = sqlite.Open(cfg.Database.DSN)
	case "mysql", "":
		dialector = mysql.Open(cfg.Database.DSN)
	default:
		return nil, fmt.Errorf("unsupported database driver %s", cfg.Database.Driver)
	}

	gormCfg := &gorm.Config{}
	if cfg.App.Env == "dev" || cfg.App.Env == "local" {
		gormCfg.Logger = gormlogger.Default.LogMode(gormlogger.Info)
	}

	db, err := gorm.Open(dialector, gormCfg)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.AutoMigrate(
		&model.User{},
		&model.AuditLog{},
		&model.ManagerEmployee{},
		&model.ContentCategory{},
		&model.Content{},
		&model.LearningRecord{},
		&model.Banner{},
		&model.ExamPaper{},
		&model.ExamQuestion{},
		&model.ExamOption{},
		&model.ExamAttempt{},
	); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}

	// 为MySQL数据库添加表注释（SQLite不支持表注释）
	if cfg.Database.Driver == "mysql" || cfg.Database.Driver == "" {
		tableComments := map[string]string{
			"users":              "用户表",
			"audit_logs":         "审计日志表",
			"manager_employees":  "店长员工关联表",
			"content_categories": "学习内容分类表",
			"contents":           "学习内容表",
			"learning_records":   "学习记录表",
			"banners":            "轮播图表",
			"exam_papers":        "试卷表",
			"exam_questions":    "考试题目表",
			"exam_options":       "考试选项表",
			"exam_attempts":      "考试记录表",
		}

		for tableName, comment := range tableComments {
			if err := db.Exec(fmt.Sprintf("ALTER TABLE `%s` COMMENT = '%s'", tableName, comment)).Error; err != nil {
				logger.Warn("failed to add table comment", zap.String("table", tableName), zap.Error(err))
			}
		}
	}

	logger.Info("database connected")
	return db, nil
}
