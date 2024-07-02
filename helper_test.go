package gormhelper_test

import (
	"context"
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
	"time"

	"github.com/nuominmin/gorm-helper"
)

type User struct {
	Id         int64  `gorm:"<-:create;column:id;primaryKey;autoIncrement"`
	Address    string `gorm:"<-:create;column:address;type:varchar(42);uniqueIndex:uni_address;comment:'地址'"`
	Nickname   string `gorm:"column:nickname;type:varchar(255);comment:'昵称'"`
	Avatar     string `gorm:"column:avatar;type:varchar(255);comment:'头像'"`
	TotalPower uint64 `gorm:"column:total_power;comment:'total power'"`
	Level      uint64 `gorm:"column:level;comment:'level'"`
}

func (User) TableName() string {
	return "user"
}

func TestGetTableName(t *testing.T) {
	tableName := gormhelper.GetTableName[User]()
	expected := "user"

	if tableName != expected {
		t.Errorf("Expected table name to be %s, but got %s", expected, tableName)
	}
}

func TestBulkCreate(t *testing.T) {
	db, err := connect()
	if err != nil {
		t.Fatalf("Failed to setup database: %v", err)
	}

	users := []*User{
		{Address: "Alice"},
		{Address: "Bob"},
		{Address: "Charlie"},
	}

	ctx := context.Background()
	err = gormhelper.BulkCreate[User](db, ctx, users, 2)
	if err != nil {
		t.Fatalf("BulkInsert failed: %v", err)
	}

	var count int64
	db.Model(&User{}).Count(&count)
	if count != 3 {
		t.Errorf("Expected 3 users, got %d", count)
	}
}

func TestSave(t *testing.T) {
	db, err := connect()
	if err != nil {
		t.Fatalf("Failed to setup database: %v", err)
	}

	user := &User{Address: "Alice"}
	ctx := context.Background()

	// Save a new user
	err = gormhelper.Save[User](db, ctx, user)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	var count int64
	db.Model(&User{}).Count(&count)
	if count != 1 {
		t.Errorf("Expected 1 user, got %d", count)
	}

	// Update the existing user
	user.Address = "Alice Updated"
	err = gormhelper.Save[User](db, ctx, user)
	if err != nil {
		t.Fatalf("Save update failed: %v", err)
	}

	var updatedUser User
	db.First(&updatedUser, user.Id)
	if updatedUser.Address != "Alice Updated" {
		t.Errorf("Expected user name to be 'Alice Updated', got '%s'", updatedUser.Address)
	}
}

func TestFirstOrCreate(t *testing.T) {
	db, err := connect()
	if err != nil {
		t.Fatalf("Failed to setup database: %v", err)
	}

	ctx := context.Background()

	// Save a new user
	data, err := gormhelper.FirstOrCreate[User](db, ctx, &User{
		Address: "Bob1",
	}, gormhelper.WithWhere("address = ?", "Bob1"))
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	t.Logf("%+v", data)
}

func TestFindWithCount(t *testing.T) {
	db, err := connect()
	if err != nil {
		t.Fatalf("Failed to setup database: %v", err)
	}

	users := []*User{
		{Address: "Alice"},
		{Address: "Bob"},
		{Address: "Charlie"},
	}
	db.Create(&users)

	ctx := context.Background()
	page, size := 1, 2
	data, total, err := gormhelper.FindWithCount[User](db, ctx, page, size, gormhelper.WithWhere("address in (?)", []string{"Bob", "Alice"}))
	if err != nil {
		t.Fatalf("FindWithCount failed: %v", err)
	}
	if total != 2 {
		t.Errorf("Expected total to be 2, got %d", total)
	}
	if len(data) != 2 {
		t.Errorf("Expected 2 users, got %d", len(data))
	}
}

func TestCount(t *testing.T) {
	db, err := connect()
	if err != nil {
		t.Fatalf("Failed to setup database: %v", err)
	}

	users := []*User{
		{Address: "Alice"},
		{Address: "Bob"},
	}
	db.Create(&users)

	ctx := context.Background()
	total, err := gormhelper.Count[User](db, ctx)
	if err != nil {
		t.Fatalf("Count failed: %v", err)
	}
	if total != 2 {
		t.Errorf("Expected total to be 2, got %d", total)
	}
}

func connect() (*gorm.DB, error) {
	conn, err := gorm.Open(mysql.Open("root:123456@tcp(127.0.0.1:3306)/hunter?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	conn.Logger = conn.Logger.LogMode(logger.LogLevel(4))

	var db *sql.DB
	if db, err = conn.DB(); err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Second * 300)
	return conn, nil
}
