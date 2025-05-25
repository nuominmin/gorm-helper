package gormhelper

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

// ConnectMysql
// dsn: root:123456@tcp(127.0.0.1:3306)/database?charset=utf8mb4&parseTime=True&loc=Local
func ConnectMysql(dsn string, opts ...ConnOption) (*gorm.DB, error) {
	options := newConnOptions(opts...)
	return Connect(mysql.Open(dsn), options)
}

// ConnectSqlite
// dsn: file:data.db?cache=shared&mode=rwc
func ConnectSqlite(dsn string, opts ...ConnOption) (*gorm.DB, error) {
	options := newConnOptions(opts...)
	return Connect(sqlite.New(sqlite.Config{
		DriverName: "sqlite",
		DSN:        dsn,
	}), options)
}

// Connect .
func Connect(dialector gorm.Dialector, options ConnOptions) (*gorm.DB, error) {
	conn, err := gorm.Open(dialector, options.config)
	if err != nil {
		return nil, fmt.Errorf("faile to open database: %v", err)
	}

	if options.logLevel > 0 {
		conn.Logger = conn.Logger.LogMode(options.logLevel)
	}

	var db *sql.DB
	if db, err = conn.DB(); err != nil {
		return nil, fmt.Errorf("failed to : %v", err)
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

	return conn, nil
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
		logLevel:        1,
		maxIdleConns:    10,
		maxOpenConns:    100,
		connMaxLifetime: time.Second * 300,
	}
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

// WithConnLogLevel default silent
// 1. silent
// 2. error
// 3. warn
// 4. info
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
