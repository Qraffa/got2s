## got2s

got2s可以将mysql数据库的表字段，mysql建表语句转换为对应的golang结构体。

### 安装

```bash
$ git clone https://github.com/Qraffa/got2s.git
$ cd got2s
$ go build -o got2s main.go
```

需要golang1.14.6以上的环境。

如果需要将生成的结构体复制到剪切板的功能，unix环境还需安装`xclip`或`xsel`

剪切板详见：[https://github.com/atotto/clipboard](https://github.com/atotto/clipboard)

### 用法

```bash
$ ./got2s help
got2s get struct from mysql table or create table statement

Usage:
  got2s [flags]
  got2s [command]

Available Commands:
  help        Help about any command
  sql         Parse create table statement to struct
  table       Convert mysql table to struct
  version     Print the version number of got2s

Flags:
  -h, --help   help for got2s

Use "got2s [command] --help" for more information about a command.
```

#### sql

将mysql建表语句转换为golang结构体，需要指定`-s`或`--sql`，表示需要转换的sql语句的文件

```bash
$ ./got2s sql -h
Parse create table statement to struct

Usage:
  got2s sql [flags]

Flags:
  -f, --func             add TableName function
  -h, --help             help for sql
  -i, --import           auto import needed package
  -j, --json             add json tag
  -p, --package string   the name of struct package
  -s, --sql string       set sql file
```

#### table

将mysql数据库的表字段转换为golang结构体，需要指定数据库的url和db和schema，表示需要转换的表

```bash
$ ./got2s table -h
Convert mysql table to struct

Usage:
  got2s table [flags]

Flags:
  -d, --db string        the database of sql
  -f, --func             add TableName function
  -h, --help             help for table
  -i, --import           auto import needed package
  -j, --json             add json tag
  -p, --package string   the name of struct package
  -t, --table string     the table of sql
  -u, --url string       the url of sql connection (default "root:123456@tcp(localhost:3306)")
```

#### 附加参数解释

- ` -f`或`--func`：表示是否添加`TableName`函数，来指定表名，默认为false
- `-i`或`--import`：表示是否需要自动导入所需要的包，默认为false
- `-j`或`--json`：表示是否需要添加json tag，默认为false
- `-p`或`--package`：表示golang结构体的包名，默认不写入包名

### example

1. 对建表语句使用

```bash
$vim create.sql
CREATE TABLE IF NOT EXISTS `blog`(
   `id` INT UNSIGNED AUTO_INCREMENT,
   `title` VARCHAR(100) NOT NULL,
   `author` VARCHAR(40) NOT NULL,
   PRIMARY KEY ( `id` )
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

$ ./got2s sql -s create.sql -fij -p model
package model

type Blog struct {
	Id     int64  `gorm:"column:id" json:"id"`
	Title  string `gorm:"column:title" json:"title"`
	Author string `gorm:"column:author" json:"author"`
}

func (*Blog) TableName() string {
	return "blog"
}
```

2. 对数据库表使用

```bash
$ ./got2s table -uroot:123456@tcp\(localhost:3306\) -d got2s_test -t blog -fij -p model
package model

type Blog struct {
	Id     int64  `gorm:"column:id" json:"id"`
	Title  string `gorm:"column:title" json:"title"`
	Author string `gorm:"column:author" json:"author"`
}

func (*Blog) TableName() string {
	return "blog"
}
```

在bash下的`()`需要使用一下转义`\`

### 实现

1. 连接数据库查询表信息

   ```sql
   SHOW FULL COLUMNS FROM table_name;
   ```

   获取各字段名，类型，注释信息

2. 使用类似语法树的结构，确定最终生成代码的结构体格式

3. 遍历查询字段的结果，确定结构体内容

4. 主要使用fmt.Sprintf来拼接各部分内容(可以试试使用`text/template`来改进一下)

#### 其他实现细节

1. sql类型转go类型

   使用go的正则匹配来转换

   有待完善...

2. 自动处理包导入

   对于生成结构体中存在`time.Time`类型的，`imports.Process`方法来处理包的自动导入

   不过一般对于时间大多数是使用时间戳，这个情况比较少

3. 格式化代码

   对于最终生成代码，使用`format.Source`方法来处理格式化

4. 解析sql语句

   [https://github.com/pingcap/parser](https://github.com/pingcap/parser)