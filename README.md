### 示例

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
