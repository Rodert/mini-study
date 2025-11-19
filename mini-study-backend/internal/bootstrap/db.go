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
	); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}

	logger.Info("database connected")
	return db, nil
}
