package gormhelper

// Model 定义一个包含 TableName 方法的接口
type Model interface {
	TableName() string
}

var (
	// MaxQuerySize 最大查询大小
	MaxQuerySize = 500
	// DefaultQuerySize 默认查询大小
	DefaultQuerySize = 20
)

func GetTableName[T Model]() string {
	var model T
	return model.TableName()
}
