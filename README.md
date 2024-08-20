# gorm-helper 是一个用于简化使用 GORM 的辅助包，提供了实际化的增删改查操作


## 方法介绍
### Create
- 在数据库中创建一条新记录。该方法支持通过 Option 参数进行定制，允许你在遇到重复数据时忽略错误

### FirstOrCreate
- 查询数据库中的第一条记录。如果未找到该记录，则使用提供的默认数据创建一条新记录

### BulkCreate
- 在数据库中批量创建多条记录。该方法可以选择性地忽略重复记录，并支持分批次插入以防止 SQL 语句过长导致的性能问题

### UpdateOrCreate
- 在数据库中根据指定条件更新数据，如果未找到对应记录，则创建新记录。它是一个事务操作，确保在并发情况下的操作原子性

### UpdateColumn
- 用于更新数据库中指定列的值。它是一种快速更新单列值的操作

### PagingParams
- 用于分页查询时的参数结构体，包含分页信息和排序信息

### Count
- 用于查询指定表的数据总数。你可以使用 Option 参数来自定义查询条件

### FindWithCount
- 用于分页查询指定表的数据，同时返回结果集和数据总数

### Find
- 用于分页查询指定表的数据并返回结果集

### FindAll
- 用于查询指定表的所有数据并返回结果集

### First
- 用于查询指定表的第一行数据。如果未找到记录且指定了 Ignore 选项，则返回 nil 而不报错

### FirstWithTransform
- 用于查询指定表的第一行数据并进行转换。你可以传入一个自定义的转换函数，将结果转换为其他数据结构

### FirstWithoutIgnore
- 用于查询指定表的第一行数据，该方法支持通过 Option 参数进行定制，允许你在遇到重复数据时忽略错误
- 不会忽略通过 Option 传入的 Ignore 设置

## 示例

``` go
type User struct {
	Id         int64  `gorm:"<-:create;column:id;primaryKey;autoIncrement"`
	Address    string `gorm:"<-:create;column:address;type:varchar(42);uniqueIndex:uni_address;comment:'地址'"`
	Nickname   string `gorm:"column:nickname;type:varchar(255);comment:'昵称'"`
	Avatar     string `gorm:"column:avatar;type:varchar(255);comment:'头像'"`
}

func (u User) TableName() string {
	return "user"
}

func (r *UserRepo) Find(ctx context.Context, page, size, level int, addresses ...string) (users []*User, total int64, err error) {
	var options []gormhelper.Option

	options = append(options,
		gormhelper.WithWhere("level <= ? and level > 2", level),
		gormhelper.WithOrderBy("total_power desc"),
	)

	if len(addresses) > 0 {
		options = append(options, gormhelper.WithWhere("address in (?)", addresses))
	}

	users, total, err = gormhelper.FindWithCount[User](r.data.client, ctx, page, size, options...)

	if err != nil {
		r.log.Errorf("find users error: %v", err)
		return nil, 0, err
	}
	return users, total, nil
}

func (r *UserRepo) LoadOnCreate(ctx context.Context, address string) (*User, error) {
	address = strings.ToLower(address)
	defaultData := &User{
		Address:  address,
		Nickname: address,
	}
	user, err := gormhelper.FirstOrCreate[User](r.data.client, ctx, defaultData,
		gormhelper.WithWhere("address = ?", address),
		gormhelper.WithIgnore(),
	)
	if err != nil {
		r.log.Errorf("FirstOrCreate user error: %v", err)
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) UpdateNickname(ctx context.Context, address, nickname string) error {
	address = strings.ToLower(address)
	defaultData := &User{
		Address:  address,
		Nickname: nickname,
	}

	updateData := map[string]interface{}{
		"nickname": nickname,
	}

	err := gormhelper.UpdateOrCreate[User](r.data.client, ctx, defaultData, updateData,
		gormhelper.WithWhere("address = ?", address),
		gormhelper.WithIgnore(),
	)

	if err != nil {
		r.log.Errorf("update avatar error: %v", err)
		return err
	}

	return nil
}

```
