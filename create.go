package gormhelper

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"strings"
)

const DUPLICATE_ENTRY = "Duplicate entry"

// Create 创建数据
func Create[T Model](db *gorm.DB, ctx context.Context, data *T, opts ...Option) error {
	err := ApplyOptions[T](db, ctx, opts...).Create(data).Error
	if err != nil {
		if options := NewOptions(opts...); options.Ignore && strings.Contains(err.Error(), DUPLICATE_ENTRY) {
			return nil
		}
		return err
	}
	return nil
}

// FirstOrCreate 查询指定表的第一行数据, 如果不存在则创建
func FirstOrCreate[T Model](db *gorm.DB, ctx context.Context, defaultData *T, opts ...Option) (data *T, err error) {
	data, err = First[T](db, ctx, opts...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return defaultData, Create(db, ctx, defaultData)
		}
		return nil, err
	}
	return data, nil
}
