package util

import (
	"fmt"
	"testing"
)

func Test_t2s(t *testing.T) {

	option := new(Option)
	option.User = "root"
	option.Password = "12345678"
	option.Host = "localhost"
	option.Port = 3306
	option.DataBase = "purple"
	option.Table = "id_bind_phone"

	option.IsSQL = true
	option.Sql = "CREATE TABLE IF NOT EXISTS `id_bind_phone` (\n\t\t`account_type` TINYINT unsigned NOT NULL DEFAULT '1' COMMENT '我是1字段的注释',\n\t\t`account_lei` TINYINT unsigned NOT NULL DEFAULT '1',\n\t\t`account_leii` TINYINT unsigned NOT NULL DEFAULT '1' COMMENT '我是3字段的注释'\n\t\t)  ENGINE=InnoDB DEFAULT CHARSET=utf8mb4"

	option.IsJson = true
	option.FormatType = LineType
	option.IsImport = true
	option.IsFunc = true
	option.IsPackage = true
	option.Package = "model"

	src, err := Parse(option)
	if err != nil {
		fmt.Println("create fail")
		fmt.Println(err.Error())
	} else {
		fmt.Println("create ok")
		fmt.Println(src)
	}
}
