package util

import (
	"database/sql"
	"fmt"
	"go/format"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/tools/imports"
	"got2s/parse"
)

// option for parse
type Option struct {
	// ***** 选择已存在表的查询参数 *****
	Host     string // 数据库host
	Port     int32  // 数据库端口
	User     string // 数据库用户名
	Password string // 数据库密码
	DataBase string // 数据库
	Table    string // 表
	Url      string // 数据库地址

	// ***** 选择sql解析参数 *****
	IsSQL bool   // 是否为sql语句查询，default false
	Sql   string // 建表sql语句

	// ***** 其他附加参数 *****
	IsJson     bool   // 是否添加json
	FormatType int32  // json字段的命名格式
	IsImport   bool   // 是否需要导入包
	IsFunc     bool   // 是否需要添加方法
	IsPackage  bool   // 是否需要包信息
	Package    string // 包名称
}

const (
	LineType  = 1 // 下划线
	LowerType = 2 // 首字母小写驼峰
	UpperType = 3 // 首字母大写驼峰
)

// CIStr is case insensitive string.
type CIStr struct {
	O string // Origin string
	F string // format string
}

// Tag 字段注解
type Tag struct {
	Key   string
	Value interface{}
}

func (t Tag) String() string {
	return fmt.Sprintf("%s:\"%s\"", t.Key, t.Value)
}

// Field 结构体字段
type Field struct {
	Name    CIStr
	Type    string
	Tags    []*Tag
	Comment string
}

func (f Field) String() string {
	tags := ""
	length := len(f.Tags)
	if length >= 1 {
		tags += f.Tags[0].String()
		for i := 1; i < length; i++ {
			tags += fmt.Sprintf(" %s", f.Tags[i])
		}
		tags = fmt.Sprintf("`%s`", tags)
	}
	// 空注释不输出
	if f.Comment == "" {
		return fmt.Sprintf("%s %s %s", f.Name.F, f.Type, tags)
	}
	return fmt.Sprintf("%s %s %s // %s", f.Name.F, f.Type, tags, f.Comment)
}

// TableStruct 表结构体
type TableStruct struct {
	Name   CIStr
	Fields []*Field
}

func (t TableStruct) String() string {
	src := ""
	for _, val := range t.Fields {
		src += fmt.Sprintf("%s\n", val)
	}
	src = fmt.Sprintf("type %s struct {\n%s}", t.Name.F, src)
	return src
}

// Code 生成代码
type Code struct {
	Package string
	Struct  *TableStruct
	Func    string
}

func (c Code) String() string {
	return fmt.Sprintf("%s\n%s\n%s\n", c.Package, c.Struct, c.Func)
}

// TableInfo 表结构信息
type TableInfo struct {
	Field      []byte
	Type       []byte
	Collation  []byte
	Null       []byte
	Key        []byte
	Default    []byte
	Extra      []byte
	Privileges []byte
	Comment    []byte
}

// sql类型转golang类型
func TransType(sqlType []byte) string {
	if match, _ := regexp.Match("int", sqlType); match {
		return "int64"
	} else if match, _ := regexp.Match("varchar", sqlType); match {
		return "string"
	} else if match, _ := regexp.Match("date", sqlType); match {
		return "time.Time"
	} else if match, _ := regexp.Match("time", sqlType); match {
		return "time.Time"
	} else if match, _ := regexp.Match("float", sqlType); match {
		return "float64"
	} else if match, _ := regexp.Match("double", sqlType); match {
		return "float64"
	} else if match, _ := regexp.Match("bool", sqlType); match {
		return "bool"
	} else {
		return "string"
	}
}

// 获取方法
func GetFunc(structName, tableName string) string {
	return fmt.Sprintf("func (*%s) TableName() string {\n return \"%s\"\n}", structName, tableName)
}

// 对于已存在表，直接查询数据库
func preQuery(cfg *Option, t *TableStruct) error {
	// username:password@protocol(address)/dbname?param=value
	//DSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DataBase)
	DSN := fmt.Sprintf("%s/%s", cfg.Url, cfg.DataBase)
	// 连接数据库
	db, err := sql.Open("mysql", DSN)
	defer db.Close()
	if err != nil {
		return err
	}
	// 查询表结构
	sql := fmt.Sprintf("SHOW FULL COLUMNS FROM %s", cfg.Table)
	rows, err := db.Query(sql)
	if err != nil {
		return err
	}
	// 遍历表结构获取字段信息
	infos := make([]*TableInfo, 0)
	for rows.Next() {
		table := new(TableInfo)
		if err := rows.Scan(&table.Field, &table.Type, &table.Collation, &table.Null, &
			table.Key, &table.Default, &table.Extra, &table.Privileges, &table.Comment); err != nil {
			return err
		}
		infos = append(infos, table)
	}
	// 遍历表字段信息，生成结构体
	t.Name.O = cfg.Table
	t.Name.F = UpperCamelCase(cfg.Table)
	for _, val := range infos {
		t.Fields = append(t.Fields, &Field{
			Name:    CIStr{string(val.Field), UpperCamelCase(string(val.Field))},
			Type:    TransType(val.Type),
			Comment: string(val.Comment),
		})
	}
	return nil
}

// 对于sql语句，直接解析
func preSQL(sql string, t *TableStruct) error {
	create := new(parse.CreateInfo)
	create, err := parse.ParseSQL(sql)
	if err != nil {
		return err
	}
	t.Name.O = create.TableName
	t.Name.F = UpperCamelCase(create.TableName)
	for k, val := range create.ColNames {
		field := new(Field)
		field.Name = CIStr{val, UpperCamelCase(val)}
		field.Type = TransType([]byte(create.ColTypes[k]))
		if comment, ok := create.ColComment[k+1]; ok {
			field.Comment = comment
		}
		t.Fields = append(t.Fields, field)
	}
	return nil
}

// 解析表
func Parse(cfg *Option) (string, error) {
	tstruct := new(TableStruct)
	// sql解析/查询表
	var err error
	if cfg.IsSQL {
		err = preSQL(cfg.Sql, tstruct)
	} else {
		err = preQuery(cfg, tstruct)
	}
	if err != nil {
		return "", err
	}
	// ******************************************************
	// 添加tag字段信息
	for _, val := range tstruct.Fields {
		tags := make([]*Tag, 0)
		ormTag := new(Tag)
		ormTag.Key = "gorm"
		ormTag.Value = fmt.Sprintf("column:%s", val.Name.O)
		tags = append(tags, ormTag)
		if cfg.IsJson {
			jsonTag := new(Tag)
			jsonTag.Key = "json"
			switch cfg.FormatType {
			case LineType:
				jsonTag.Value = LineCase(val.Name.O)
			case LowerType:
				jsonTag.Value = LowerCamelCase(val.Name.O)
			case UpperType:
				jsonTag.Value = UpperCamelCase(val.Name.O)
			default:
				jsonTag.Value = LineCase(val.Name.O)
			}
			tags = append(tags, jsonTag)
		}
		val.Tags = tags
	}
	// 添加包信息
	pack := ""
	if cfg.IsPackage {
		pack = fmt.Sprintf("package %s", cfg.Package)
	}
	// 添加方法
	funcSrc := ""
	if cfg.IsFunc {
		funcSrc = fmt.Sprintf("%s", GetFunc(tstruct.Name.F, tstruct.Name.O))
	}
	code := Code{Package: pack, Struct: tstruct, Func: funcSrc}
	src := code.String()
	// 添加导入包
	if cfg.IsImport {
		options := &imports.Options{
			TabWidth:   8,
			TabIndent:  true,
			Comments:   true,
			Fragment:   true,
			FormatOnly: false,
		}
		res, err := imports.Process("", []byte(src), options)
		if err != nil {
			return "", err
		}
		src = string(res)
	}
	// 结构体格式化
	res, err := format.Source([]byte(src))
	if err != nil {
		return "", err
	}
	src = string(res)
	return src, nil
}
