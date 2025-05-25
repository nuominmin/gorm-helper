package gormhelper_test

import (
	"context"
	"github.com/nuominmin/gorm-helper"
	"testing"
)

const (
	mysqlDsn  = "root:123456@tcp(127.0.0.1:3306)/hunter?charset=utf8mb4&parseTime=True&loc=Local"
	sqliteDsn = "file:data.db?cache=shared&mode=rwc"
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
	db, err := gormhelper.ConnectMysql(mysqlDsn)
	if err != nil {
		t.Fatalf("Failed to setup database: %v", err)
	}
	db.AutoMigrate(&User{})

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

func TestFirstOrCreate(t *testing.T) {
	db, err := gormhelper.ConnectMysql(mysqlDsn)
	if err != nil {
		t.Fatalf("Failed to setup database: %v", err)
	}
	db.AutoMigrate(&User{})

	ctx := context.Background()

	// Save a new user
	data, err := gormhelper.FirstOrCreate[User](db, ctx, &User{
		Address: "Bo1b1",
	}, gormhelper.WithWhere("address = ?", "Bob1"))
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	t.Logf("%+v", data)
}

func TestUpdateOrCreate(t *testing.T) {
	db, err := gormhelper.ConnectMysql(mysqlDsn)
	if err != nil {
		t.Fatalf("Failed to setup database: %v", err)
	}
	db.AutoMigrate(&User{})

	ctx := context.Background()

	defaultData := &User{
		Address: "Bob1",
	}

	updateData := map[string]interface{}{
		"address": "xxxxxx222xxx1xx",
		"level":   "2",
	}

	// Save a new user
	var data *User
	data, err = gormhelper.Upsert[User](db, ctx, defaultData, updateData, gormhelper.WithWhere("address = ?", "xxxxxx222xxx1xx"))
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	t.Logf("%+v", data)
}

func TestFindWithCount(t *testing.T) {
	db, err := gormhelper.ConnectMysql(mysqlDsn)
	if err != nil {
		t.Fatalf("Failed to setup database: %v", err)
	}
	db.AutoMigrate(&User{})

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
	db, err := gormhelper.ConnectMysql(mysqlDsn)
	if err != nil {
		t.Fatalf("Failed to setup database: %v", err)
	}
	db.AutoMigrate(&User{})

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
