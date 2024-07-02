package gormhelper

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"reflect"
	"strings"
)

// BulkCreate 批量创建
func BulkCreate[T Model](db *gorm.DB, ctx context.Context, data []*T, batchSize int, opts ...Option) (err error) {
	if len(data) == 0 {
		return nil
	}

	options := NewOptions(opts...)
	tableName := GetTableName[T]()
	columns, placeholders := getColumnsAndPlaceholders(data[0])

	for i := 0; i < len(data); i += batchSize {
		end := i + batchSize
		if end > len(data) {
			end = len(data)
		}

		sql := "INSERT "
		if options.Ignore {
			sql += "IGNORE "
		}
		sql += "INTO " + tableName + " (" + columns + ") VALUES "
		var values []interface{}
		for _, record := range data[i:end] {
			sql += "(" + placeholders + "),"
			values = append(values, getValues(record)...)
		}
		sql = sql[:len(sql)-1] // 去掉最后一个逗号

		if err = db.WithContext(ctx).Exec(sql, values...).Error; err != nil {
			return err
		}
	}

	return nil
}

// 获取数据的列名和占位符
func getColumnsAndPlaceholders[T Model](data *T) (string, string) {
	val := reflect.ValueOf(data).Elem()
	numFields := val.NumField()
	columns := make([]string, numFields)
	placeholders := make([]string, numFields)
	for i := 0; i < numFields; i++ {
		placeholders[i] = "?"

		field := val.Type().Field(i)

		if tag, ok := field.Tag.Lookup("gorm"); ok {
			tagSettings := schema.ParseTagSetting(tag, ";")
			var column string
			if column, ok = tagSettings["COLUMN"]; ok {
				columns[i] = column
				continue
			}
		}

		columns[i] = field.Name
	}
	return strings.Join(columns, ", "), strings.Join(placeholders, ", ")
}

// 获取数据的值
func getValues[T Model](data *T) []interface{} {
	val := reflect.ValueOf(data).Elem()
	values := make([]interface{}, val.NumField())
	for i := 0; i < val.NumField(); i++ {
		values[i] = val.Field(i).Interface()
	}
	return values
}
