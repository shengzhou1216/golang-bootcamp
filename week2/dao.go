package week2

import (
	"database/sql"
	"errors"
	"fmt"
)

// 方法1: 不认为sql.ErrNoRows是错误
func query1() (int, error) {
	err := sql.ErrNoRows
	if errors.Is(sql.ErrNoRows,err) {
		return 0,nil
	}
	return 0, err
}

// 方法2: 包装错误, 向上抛
func query2() (int, error) {
	err := sql.ErrNoRows
	if errors.Is(sql.ErrNoRows,err) {
		return 0,fmt.Errorf("no more data: %w",err)
	}
	return 0, err
}



