package gormhelper

import (
	"context"
	"gorm.io/gorm"
)

// Count 查询指定表的数据总数
func Count[T Model](db *gorm.DB, ctx context.Context, opts ...Option) (total int64, err error) {
	return total, ApplyOptions[T](db, ctx, opts...).Count(&total).Error
}

// FindWithCount 查询指定表的数据并返回结果集以及总数
func FindWithCount[T Model](db *gorm.DB, ctx context.Context, page, size int, opts ...Option) (data []*T, total int64, err error) {
	// 计算总数
	if total, err = Count[T](db, ctx, opts...); err != nil {
		return nil, 0, err
	}

	offset, limit := PagingParams(page, size)
	return data, total, ApplyOptions[T](db, ctx, opts...).Offset(offset).Limit(limit).Find(&data).Error
}

// FindAll 查询所有
func FindAll[T Model](db *gorm.DB, ctx context.Context, opts ...Option) (data []*T, err error) {
	return data, ApplyOptions[T](db, ctx, opts...).Find(&data).Error
}

// First 查询指定表的第一行数据
func First[T Model](db *gorm.DB, ctx context.Context, opts ...Option) (data *T, err error) {
	return data, ApplyOptions[T](db, ctx, opts...).First(&data).Error
}
