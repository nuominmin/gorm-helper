package gormhelper

import (
	"context"
	"gorm.io/gorm"
)

type Options struct {
	OrderBy string
	Wheres  []WhereClause
	Ignore  bool
}

type WhereClause struct {
	Query interface{}
	Args  []interface{}
}

type Option func(*Options)

func NewOptions(opts ...Option) Options {
	var options Options
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

// ApplyOptions 应用选项到查询中
func ApplyOptions[T Model](db *gorm.DB, ctx context.Context, opts ...Option) *gorm.DB {
	options := NewOptions(opts...)
	query := db.WithContext(ctx).Table(GetTableName[T]())
	for i := 0; i < len(options.Wheres); i++ {
		query = query.Where(options.Wheres[i].Query, options.Wheres[i].Args...)
	}
	if options.OrderBy != "" {
		query = query.Order(options.OrderBy)
	}
	return query
}

// WithOrderBy 用于设置排序字段
func WithOrderBy(orderBy string) Option {
	return func(opts *Options) {
		opts.OrderBy = orderBy
	}
}

// WithWhere 用于设置 where 条件
func WithWhere(query interface{}, args ...interface{}) Option {
	return func(opts *Options) {
		opts.Wheres = append(opts.Wheres, WhereClause{Query: query, Args: args})
	}
}

func WithIgnore() Option {
	return func(opts *Options) {
		opts.Ignore = true
	}
}
