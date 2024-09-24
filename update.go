package gormhelper

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

// Upsert 根据 where 检索数据更新，如果不存在则创建数据
func Upsert[T Model](db *gorm.DB, ctx context.Context, defaultData *T, updateData map[string]interface{}, opts ...Option) (*T, error) {
	if defaultData == nil {
		return nil, errors.New("defaultData is nil")
	}
	if updateData == nil {
		return nil, errors.New("updateData is nil")
	}
	if len(opts) == 0 {
		return nil, errors.New("opts is nil")
	}
	if options := NewOptions(opts...); len(options.Wheres) == 0 {
		return nil, errors.New("wheres is nil")
	}

	var result T
	return &result, db.Transaction(func(tx *gorm.DB) error {
		// 获取模型的主键字段名称
		stmt := &gorm.Statement{DB: tx}
		if err := stmt.Parse(&result); err != nil {
			return err
		}
		if len(stmt.Schema.PrimaryFields) == 0 {
			return errors.New("no primary key field found")
		}
		primaryKeyFieldName := stmt.Schema.PrimaryFields[0].Name
		primaryKeyFieldDBName := stmt.Schema.PrimaryFields[0].DBName

		if err := ApplyOptions[T](tx, ctx, opts...).Select(primaryKeyFieldDBName).First(&result).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 创建新数据
				if err = ApplyOptions[T](tx, ctx, opts...).Create(defaultData).Error; err != nil {
					if strings.Contains(err.Error(), DUPLICATE_ENTRY) {
						// 更新成功，获取更新后的数据
						if err = Apply[T](tx, ctx).Where(defaultData).First(&result).Error; err != nil {
							return fmt.Errorf("get updated data failed, error: %v", err)
						}
						return nil
					}
					return fmt.Errorf("create data failed, error: %v", err)
				}
				result = *defaultData
				return nil
			}
			return fmt.Errorf("find data failed, error: %v", err)
		}

		primaryKeyValue := reflect.ValueOf(result).FieldByName(primaryKeyFieldName).Interface()

		// 如果找到了记录，则进行更新
		if err := Apply[T](tx, ctx).Where(primaryKeyFieldDBName+" = ?", primaryKeyValue).Updates(updateData).Error; err != nil {
			return fmt.Errorf("data found but update failed, error: %v", err)
		}

		// 更新成功，获取更新后的数据
		if err := tx.First(&result, primaryKeyValue).Error; err != nil {
			return fmt.Errorf("get updated data failed, error: %v", err)
		}
		return nil
	})
}

// 判断是否唯一键冲突错误
func IsDuplicateEntryError(err error) bool {
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}

// UpdateColumn 更新单列
func UpdateColumn[T Model](db *gorm.DB, ctx context.Context, column string, value interface{}, opts ...Option) (err error) {
	return ApplyOptions[T](db, ctx, opts...).Update(column, value).Error
}
