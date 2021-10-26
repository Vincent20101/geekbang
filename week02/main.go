package main

import (
	"database/sql"
	"fmt"

	"geekbang/week02/dao"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

func main(){
	var (
		user = "root"
		password = "root"
		host = "127.0.0.1"
		port = "3306"
		dbName = "geekbang"
	)
	sqlConStr := fmt.Sprintf("%s:%s@(%s:%s)/%s", user, password, host, port, dbName)
	var err error
	dao.DBConn, err = sql.Open("mysql", sqlConStr)
	if err != nil {
		fmt.Printf("originnal error: %T %v\n", errors.Cause(err), errors.Cause(err))
		fmt.Printf("stack trace: \n%+v", err)
		return
	}
	defer dao.DBConn.Close()

	// 业务逻辑查询
	userInfo, err := (&dao.User{}).GetUser(1)

	if err != nil {
		fmt.Printf("originnal error: %T %v\n", errors.Cause(err), errors.Cause(err))
		fmt.Printf("stack trace: \n%+v", err)
		return
	}

	fmt.Println(userInfo)
}