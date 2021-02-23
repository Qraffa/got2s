package parse

import (
	"fmt"
	"testing"
)

func Test_Parse(t *testing.T) {
	//sql := "CREATE TABLE IF NOT EXISTS `xhx`.`id_bind_phone` (" +
	//	"`id` INT(11) NOT NULL," +
	//	"`account_type` TINYINT unsigned NOT NULL DEFAULT '0' COMMENT '对应account的Type'," +
	//	"`phone` VARCHAR(20) NOT NULL DEFAULT ''," +
	//	"`last_phone` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '之前绑定的手机号phone,phone'," +
	//	"`update_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新手机号绑定时间'," +
	//	"`create_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '初绑时间'," +
	//	"PRIMARY KEY (`id`)," +
	//	"INDEX `phone` (`phone`)" +
	//	")  ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='id绑定手机号';"
	sql := "CREATE TABLE IF NOT EXISTS `xhx`.`id_bind_phone` (" +
		"`account_type` TINYINT unsigned NOT NULL DEFAULT '1' COMMENT '我是1字段的注释'," +
		"`account_lei` TINYINT unsigned NOT NULL DEFAULT '1'," +
		"`account_leii` TINYINT unsigned NOT NULL DEFAULT '1' COMMENT '我是3字段的注释'" +
		")  ENGINE=InnoDB DEFAULT CHARSET=utf8mb4"
	fmt.Println(ParseSQL(sql))
	create := Run(sql)
	fmt.Println(create.String())
	fmt.Println(create.ColComment)
}

func Test_Int(t *testing.T) {
}
