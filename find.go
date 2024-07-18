package gormhelper

import (
	"context"
	"errors"
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

	if total == 0 {
		return make([]*T, 0), 0, nil
	}

	if data, err = Find[T](db, ctx, page, size, opts...); err != nil {
		return nil, 0, err
	}
	return data, total, nil
}

// Find 查询指定表的数据并返回结果集
func Find[T Model](db *gorm.DB, ctx context.Context, page, size int, opts ...Option) (data []*T, err error) {
	offset, limit := PagingParams(page, size)
	return data, ApplyOptions[T](db, ctx, opts...).Offset(offset).Limit(limit).Find(&data).Error
}

// FindAll 查询所有
func FindAll[T Model](db *gorm.DB, ctx context.Context, opts ...Option) (data []*T, err error) {
	return data, ApplyOptions[T](db, ctx, opts...).Find(&data).Error
}

// First 查询指定表的第一行数据
func First[T Model](db *gorm.DB, ctx context.Context, opts ...Option) (data *T, err error) {
	if err = ApplyOptions[T](db, ctx, opts...).First(&data).Error; err != nil {
		if options := NewOptions(opts...); options.Ignore && errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return data, err
}

// FirstWithTransform 查询指定表的第一行数据并进行转换
func FirstWithTransform[T Model, D any](db *gorm.DB, ctx context.Context, transformFunc func(*T) *D, opts ...Option) (dto *D, err error) {
	var data *T
	if err = ApplyOptions[T](db, ctx, opts...).First(&data).Error; err != nil {
		if options := NewOptions(opts...); options.Ignore && errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return transformFunc(data), nil
}

// FirstWithoutIgnore 查询指定表的第一行数据 (不忽略 opts 传入的 Ignore)
func FirstWithoutIgnore[T Model](db *gorm.DB, ctx context.Context, opts ...Option) (data *T, err error) {
	return data, ApplyOptions[T](db, ctx, opts...).First(&data).Error
}
