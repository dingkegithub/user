package dao

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func MysqlInit(host string, port string, user, pwd, dbName string) error {
	var err error
	db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", user, pwd, host, port, dbName))
	if err != nil {
		return err
	}

	db.SingularTable(true)

	return nil
}
