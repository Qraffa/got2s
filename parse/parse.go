package parse

import (
	"fmt"
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/test_driver"
	_ "github.com/pingcap/parser/test_driver"
)

type CreateInfo struct{
	TableName string
	ColNames []string
	ColTypes []string
	ColComment map[int]string
}

func (v *CreateInfo) String() string {
	return fmt.Sprintf("表名 %s\n字段 %s\n类型 %s\n", v.TableName, v.ColNames, v.ColTypes)
}

// 回调函数 Enter
func (v *CreateInfo) Enter(in ast.Node) (ast.Node, bool) {
	// 获取注释信息
	if val, ok := in.(*ast.ColumnOption); ok {
		// 只关注注释类型的信息，其他类型跳过
		if val.Tp == ast.ColumnOptionComment {
			return in, false
		} else {
			return in, true
		}
	}
	// 获取注释内容
	if val, ok := in.(*test_driver.ValueExpr); ok {
		v.ColComment[len(v.ColNames)] = val.Datum.GetString()
	}
	// 获取字段信息
	if name, ok := in.(*ast.ColumnDef); ok {
		// 获取字段名
		v.ColNames = append(v.ColNames, name.Name.Name.O)
		// 获取字段类型
		v.ColTypes = append(v.ColTypes, name.Tp.String())
	}
	// 表名
	if name, ok := in.(*ast.TableName); ok {
		v.TableName = name.Name.O
	}
	return in, false
}

func (v *CreateInfo) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func parse(sql string) (*ast.StmtNode, error) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return nil, err
	}
	return &stmtNodes[0], nil
}

func extract(rootNode *ast.StmtNode) *CreateInfo {
	v := &CreateInfo{}
	v.ColComment = make(map[int]string)
	(*rootNode).Accept(v)
	return v
}

// 解析sql语句，返回TableStruct
func ParseSQL(sql string) (*CreateInfo, error) {
	astNode, err := parse(sql)
	if err != nil {
		return nil, err
	}
	create := extract(astNode)
	return create, nil
}

func Run(sql string) CreateInfo {
	astNode, err := parse(sql)
	if err != nil {
		fmt.Printf("parse error: %v\n", err.Error())
		return CreateInfo{}
	}
	v := extract(astNode)
	return *v
}