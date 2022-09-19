package schema

import (
	"github.com/869413421/orm/dialect"
	"go/ast"
	"reflect"
)

// Field 字段详细信息
type Field struct {
	Name string // 字段名称
	Type string // 字段类型
	Tag  string // 约束
}

type Schema struct {
	Model      interface{}       // 模型
	TableName  string            // 表名
	Fields     []*Field          // 所有字段详细信息
	FieldNames []string          // 所有字段名
	fieldMap   map[string]*Field // 字段名对应的field
}

// GetField 获取字段详细信息
func (s *Schema) GetField(name string) *Field {
	return s.fieldMap[name]
}

// Parse 获取结构体的详细信息
func Parse(dest interface{}, d dialect.Dialect) *Schema {
	// 1.获取结构体实例
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type() // TypeOf() 和 ValueOf() 是 reflect 包最为基本也是最重要的 2 个方法，分别用来返回入参的类型和值。因为设计的入参是一个对象的指针，因此需要 reflect.Indirect() 获取指针指向的实例。
	schema := &Schema{
		Model:     dest,
		TableName: modelType.Name(),
		fieldMap:  make(map[string]*Field),
	}

	// 2.循环遍历结构体的字段
	for i := 0; i < modelType.NumField(); i++ {
		// 2.1 如果是有效字段，添加到详细信息中
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup("orm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}

	return schema
}
