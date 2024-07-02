package gormhelper

// PagingParams 分页参数
func PagingParams(page, size int) (offset, limit int) {
	// 设置默认值, 如果未设置, 最大获取数据将会是 DefaultSize
	if page <= 0 {
		page = 1
	}
	if size == 0 || size > MaxQuerySize {
		size = DefaultQuerySize
	}

	// 设置最大限制, 大于 MaxSize, 将会把 size 重新设置为 DefaultSize
	return (page - 1) * size, size
}
