package schema

import (
	"geeorm/dialect"
	"go/ast"
	"reflect"
)

//一个字段
type Field struct {
	Name string
	Type string
	Tag  string
}

//存储对象与数据库数据的映射关系
type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Field
	Field_map  map[string]*Field
	Fieldnames []string
}

func (s *Schema) GetField(name string) *Field {
	return s.Field_map[name]
}


//解析，生成schema
func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:     dest,
		Name:      modelType.Name(),
		Field_map: make(map[string]*Field),
	}
	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup("geeorm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.Fieldnames = append(schema.Fieldnames, p.Name)
			schema.Field_map[p.Name] = field
		}
	}
	return schema
}
