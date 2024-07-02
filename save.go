package gormhelper

import (
	"context"
	"gorm.io/gorm"
)

// Save 保存数据，根据 where 检索数据更新，如果不存在则创建数据
func Save[T Model](db *gorm.DB, ctx context.Context, data *T, opts ...Option) (err error) {
	var total int64
	if total, err = Count[T](db, ctx, opts...); err != nil {
		return err
	}

	if total > 0 {
		// 更新数据
		return ApplyOptions[T](db, ctx, opts...).Updates(data).Error
	}

	return db.WithContext(ctx).Table(GetTableName[T]()).Create(data).Error
}
