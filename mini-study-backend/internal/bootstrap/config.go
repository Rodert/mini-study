package bootstrap

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config holds the global application configuration loaded via Viper.
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Upload   UploadConfig   `mapstructure:"upload"`
	Swagger  SwaggerConfig  `mapstructure:"swagger"`
}

// AppConfig describes metadata for the running service.
type AppConfig struct {
	Name     string `mapstructure:"name"`
	Env      string `mapstructure:"env"`
	Version  string `mapstructure:"version"`
	LogLevel string `mapstructure:"log_level"`
}

// ServerConfig includes HTTP server options.
type ServerConfig struct {
	Host              string        `mapstructure:"host"`
	Port              int           `mapstructure:"port"`
	RequestTimeoutRaw string        `mapstructure:"request_timeout"`
	ReadTimeoutRaw    string        `mapstructure:"read_timeout"`
	WriteTimeoutRaw   string        `mapstructure:"write_timeout"`
	AllowedOrigins    []string      `mapstructure:"allowed_origins"`
	RequestTimeout    time.Duration `mapstructure:"-"`
	ReadTimeout       time.Duration `mapstructure:"-"`
	WriteTimeout      time.Duration `mapstructure:"-"`
}

// DatabaseConfig holds the persistence configuration.
type DatabaseConfig struct {
	Driver string `mapstructure:"driver"`
	DSN    string `mapstructure:"dsn"`
}

// JWTConfig controls token generation.
type JWTConfig struct {
	Issuer        string        `mapstructure:"issuer"`
	Secret        string        `mapstructure:"secret"`
	TTLRaw        string        `mapstructure:"ttl"`
	RefreshTTLRaw string        `mapstructure:"refresh_ttl"`
	TTL           time.Duration `mapstructure:"-"`
	RefreshTTL    time.Duration `mapstructure:"-"`
}

// UploadConfig stores file upload limits.
type UploadConfig struct {
	MaxSizeMB int    `mapstructure:"max_size_mb"`
	Dir       string `mapstructure:"dir"`
}

// SwaggerConfig toggles Swagger UI exposure.
type SwaggerConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

// LoadConfig loads the base config plus environment overrides.
func LoadConfig(configDir string) (*Config, error) {
	v := viper.New()
	v.AddConfigPath(configDir)
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = v.GetString("app.env")
	}

	if env != "" {
		v.SetConfigName(fmt.Sprintf("config.%s", env))
		if err := v.MergeInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return nil, fmt.Errorf("merge %s config: %w", env, err)
			}
		}
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if err := cfg.normalize(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) normalize() error {
	var err error

	c.Server.RequestTimeout, err = time.ParseDuration(defaultString(c.Server.RequestTimeoutRaw, "30s"))
	if err != nil {
		return fmt.Errorf("parse server.request_timeout: %w", err)
	}

	c.Server.ReadTimeout, err = time.ParseDuration(defaultString(c.Server.ReadTimeoutRaw, "15s"))
	if err != nil {
		return fmt.Errorf("parse server.read_timeout: %w", err)
	}

	c.Server.WriteTimeout, err = time.ParseDuration(defaultString(c.Server.WriteTimeoutRaw, "15s"))
	if err != nil {
		return fmt.Errorf("parse server.write_timeout: %w", err)
	}

	c.JWT.TTL, err = time.ParseDuration(defaultString(c.JWT.TTLRaw, "24h"))
	if err != nil {
		return fmt.Errorf("parse jwt.ttl: %w", err)
	}

	c.JWT.RefreshTTL, err = time.ParseDuration(defaultString(c.JWT.RefreshTTLRaw, "168h"))
	if err != nil {
		return fmt.Errorf("parse jwt.refresh_ttl: %w", err)
	}

	if c.App.Env == "" {
		c.App.Env = "local"
	}

	return nil
}

func defaultString(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
