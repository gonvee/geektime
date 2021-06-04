package dao

import (
	"database/sql"
	"errors"
	"log"
)

var (
	db                *sql.DB
	ErrRecordNotFound = errors.New("no record found")
)

type User struct {
	Id     string
	Name   string
	Gender int
}

func init() {
	sqlDb, err := sql.Open("mysql", "**dsn**")
	if err != nil {
		log.Fatalf("dao: open mysql failed: %s", err)
	}
	db = sqlDb
}

// UserById 根据用户ID查询
// 当未查询到结果时返回ErrRecordNotFound
//
// 应该对sql.ErrNoRows做一个转换而非包装, 转换为dao层定义的error, 原因如下:
// 1. 在未找到数据时http一般返回404, 此时代码中要判断相应的error类型;
// 而对sql.ErrNoRows进行转换是为了调用dao层的包只需要依赖dao层,
// 不需要直接导入sql包。除ErrRecordNotFound外的error可以认为是内
// 部错误, 统一返回http 500状态码, 此时dao层的caller不需要具体的类
// 型信息。
// 2. wrap error更多的是为了获取详细的上下文信息, 对于sql.ErrNoRows
// 这种表意很明确的error, 调用方已经有了入参, 在业务处理中加入日志, 配
// 合dao层相关的sql日志, 用于定位分析的信息已经很充分。
func UserById(id string) (*User, error) {
	user := &User{}

	row := db.QueryRow("select id, name, gender from users where id = ?", id)

	err := row.Scan(&user.Id, &user.Name, &user.Gender)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return user, err
}
