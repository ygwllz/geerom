package dialect

import "reflect"

type Dialect interface {
	DataTypeOf(typ reflect.Value) string
	TableExistSQL(tableName string) (string, []interface{})
}

var dialectmap = map[string]Dialect{}

func RegisterDialect(name string, dialect Dialect) {
	dialectmap[name] = dialect
}

func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectmap[name]
	return
}


