package gormhelper

import (
	"context"
	"gorm.io/gorm"
)

// UpdateOrCreate 根据 where 检索数据更新，如果不存在则创建数据
func UpdateOrCreate[T Model](db *gorm.DB, ctx context.Context, defaultData *T, updateData map[string]interface{}, opts ...Option) (err error) {
	tx := db.Begin()

	var total int64
	if total, err = Count[T](tx, ctx, opts...); err != nil {
		tx.Rollback()
		return err
	}

	if total > 0 {
		// 更新数据
		if err = ApplyOptions[T](tx, ctx, opts...).Updates(updateData).Error; err != nil {
			tx.Rollback()
			return err
		}
		return nil
	}

	if err = Create(tx, ctx, defaultData); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
