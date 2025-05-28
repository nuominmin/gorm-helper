package gormhelper

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectMysql
// dsn: root:123456@tcp(127.0.0.1:3306)/database?charset=utf8mb4&parseTime=True&loc=Local
func ConnectMysql(dsn string, opts ...ConnOption) (*gorm.DB, func(), error) {
	options := newConnOptions(opts...)
	return Connect(mysql.Open(dsn), options)
}

// ConnectSqlite
// dsn: file:data.db?cache=shared&mode=rwc
func ConnectSqlite(dsn string, opts ...ConnOption) (*gorm.DB, func(), error) {
	options := newConnOptions(opts...)
	return Connect(sqlite.New(sqlite.Config{
		DriverName: "sqlite",
		DSN:        dsn,
	}), options)
}

// Connect .
func Connect(dialector gorm.Dialector, options ConnOptions) (*gorm.DB, func(), error) {
	conn, err := gorm.Open(dialector, options.config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open database: %w", err)
	}

	if options.logLevel != logger.Silent {
		conn.Logger = conn.Logger.LogMode(options.logLevel)
	}

	var db *sql.DB
	if db, err = conn.DB(); err != nil {
		return nil, nil, fmt.Errorf("failed to get underlying sql.DB: %v", err)
	}

	if options.maxIdleConns > 0 {
		db.SetMaxIdleConns(options.maxIdleConns)
	}

	if options.maxOpenConns > 0 {
		db.SetMaxOpenConns(options.maxOpenConns)
	}

	if options.connMaxLifetime > 0 {
		db.SetConnMaxLifetime(options.connMaxLifetime)
	}
	cleanup := func() {
		_ = db.Close()
	}
	return conn, cleanup, nil
}

type ConnOptions struct {
	config          *gorm.Config
	logLevel        logger.LogLevel
	maxIdleConns    int
	maxOpenConns    int
	connMaxLifetime time.Duration
}

type ConnOption func(*ConnOptions)

func newConnOptions(opts ...ConnOption) ConnOptions {
	options := ConnOptions{
		config:          &gorm.Config{},
		logLevel:        logger.Silent,
		maxIdleConns:    10,
		maxOpenConns:    100,
		connMaxLifetime: time.Second * 300,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(&options)
		}
	}
	return options
}

// WithConnLogLevel default silent
// 可选值: logger.Silent, logger.Error, logger.Warn, logger.Info
func WithConnLogLevel(logLevel logger.LogLevel) ConnOption {
	return func(opts *ConnOptions) {
		opts.logLevel = logLevel
	}
}

// WithMaxIdleConns default 10
// sets the maximum number of connections in the idle
func WithMaxIdleConns(maxIdleConns int) ConnOption {
	return func(opts *ConnOptions) {
		opts.maxIdleConns = maxIdleConns
	}
}

// WithMaxOpenConns default 100
// sets the maximum number of open connections to the database
func WithMaxOpenConns(maxOpenConns int) ConnOption {
	return func(opts *ConnOptions) {
		opts.maxOpenConns = maxOpenConns
	}
}

// WithConnMaxLifetime default time.Second * 300
// sets the maximum amount of time a connection may be reused
func WithConnMaxLifetime(connMaxLifetime time.Duration) ConnOption {
	return func(opts *ConnOptions) {
		opts.connMaxLifetime = connMaxLifetime
	}
}

// WithConnOpenConfig gorm.Config
func WithConnOpenConfig(c *gorm.Config) ConnOption {
	return func(opts *ConnOptions) {
		opts.config = c
	}
}

// GetDBStats 获取数据库连接池统计信息
func GetDBStats(db *gorm.DB) (sql.DBStats, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return sql.DBStats{}, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	return sqlDB.Stats(), nil
}

// HealthCheck 检查数据库连接健康状态
func HealthCheck(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return sqlDB.PingContext(ctx)
}
