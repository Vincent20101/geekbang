package dao

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)
var DBConn *sql.DB


type User struct {
	ID int
	Name string
}

func (u *User) GetUser(id int) (*User, error) {
	fmt.Println(id)
	user := User{}
	err := DBConn.QueryRow("SELECT id, name FROM users WHERE id in=(?)", id).Scan(&user.ID, &user.Name)
	if err != nil {
		if errors.Is(err,sql.ErrNoRows) {
			// sql.ErrNoRows 如果上层不关心数据是否为空的场景
			return &user, nil
			// 如果是第三方调用或上层关注数据是否为空的场景,返回如下
			//return nil, errors.Wrap(err, "userDao: sql.ErrNoRows")
		} else {
			// 其他错误， Warp 后交给上一层，或者记录一下日志然后返回空 + 错误
			return nil, errors.Wrap(err, "userDao: get users failed")
		}
	}
	return &user,nil
}