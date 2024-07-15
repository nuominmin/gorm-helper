package gormhelper

import (
	"context"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"reflect"
)

// UpdateOrCreate 根据 where 检索数据更新，如果不存在则创建数据
func UpdateOrCreate[T Model](db *gorm.DB, ctx context.Context, defaultData *T, updateData map[string]interface{}, opts ...Option) (err error) {
	return db.Transaction(func(tx *gorm.DB) error {
		var result *T
		// 获取模型的主键字段名称
		stmt := &gorm.Statement{DB: tx}
		if err = stmt.Parse(&result); err != nil {
			return err
		}
		if len(stmt.Schema.PrimaryFields) == 0 {
			return errors.New("no primary key field found")
		}
		primaryKeyFieldName := stmt.Schema.PrimaryFields[0].Name
		primaryKeyFieldDBName := stmt.Schema.PrimaryFields[0].DBName

		if err = ApplyOptions[T](db, ctx, opts...).Select(primaryKeyFieldDBName).First(&result).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 创建新数据
				return Create(tx, ctx, defaultData)
			}
			return err
		}

		// 将 result 转为 map
		var resultBytes []byte
		if resultBytes, err = json.Marshal(result); err != nil {
			return err
		}
		resultMap := make(map[string]interface{})
		if err = json.Unmarshal(resultBytes, &resultMap); err != nil {
			return err
		}

		// 从 map 中获取主键字段的值
		primaryKeyValue, ok := resultMap[primaryKeyFieldName]
		if !ok {
			return errors.New("primary key field value not found")
		}

		// 如果找到了记录，则进行更新
		return Apply[T](tx, ctx).Where(primaryKeyFieldDBName+" = ?", primaryKeyValue).Updates(updateData).Error
	})
}

// 通过字段名获取字段值
func getFieldValueByName(model interface{}, fieldName string) (interface{}, error) {
	v := reflect.ValueOf(model).Elem()
	f := v.FieldByName(fieldName)
	if !f.IsValid() {
		return nil, errors.New("field not found")
	}
	return f.Interface(), nil
}
