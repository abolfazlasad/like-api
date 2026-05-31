package postgres

import (
	"fmt"
	"like-api/internal/infrastructure/database/postgres/models"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB holds separate write and read database connections.
// Currently, both connections point to the same database instance.
// When a read replica becomes available, only ReadDB needs to be changed.
type DB struct {
	WriteDB *gorm.DB
	ReadDB  *gorm.DB
}

// Config contains PostgreSQL connection settings.
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// buildDSN builds a PostgreSQL DSN string from the provided configuration.
func buildDSN(cfg Config) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)
}

// newGormDB creates and configures a GORM database connection.
func newGormDB(dsn string, readOnly bool) (*gorm.DB, error) {
	logLevel := logger.Info
	if readOnly {
		// Read replicas typically generate fewer logs.
		logLevel = logger.Warn
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Connection pool configuration.
	// Read replicas usually have larger pools because they mainly serve SELECT queries.
	if readOnly {
		sqlDB.SetMaxOpenConns(20)
		sqlDB.SetMaxIdleConns(10)
	} else {
		sqlDB.SetMaxOpenConns(10)
		sqlDB.SetMaxIdleConns(5)
	}

	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// NewDB creates both write and read database connections.
// At the moment, both use the same host.
// When a read replica is introduced, create a separate read DSN.
func NewDB(cfg Config) (*DB, error) {
	dsn := buildDSN(cfg)

	writeDB, err := newGormDB(dsn, false)
	if err != nil {
		return nil, fmt.Errorf("write db: %w", err)
	}

	// TODO: Use cfg.ReadHost (or a dedicated read configuration)
	// when a read replica becomes available.
	//
	// Example:
	// readDSN := buildDSN(Config{
	//     Host: cfg.ReadHost,
	//     ...
	// })
	readDB, err := newGormDB(dsn, true)
	if err != nil {
		return nil, fmt.Errorf("read db: %w", err)
	}

	log.Println("[postgres] write and read connections established")

	return &DB{
		WriteDB: writeDB,
		ReadDB:  readDB,
	}, nil
}

// Migrate runs auto-migrations for all database models.
func (d *DB) Migrate() error {
	log.Println("[postgres] running migrations...")

	return d.WriteDB.AutoMigrate(
		&models.UserModel{},
		&models.PostModel{},
		&models.LikeModel{},
	)
}
