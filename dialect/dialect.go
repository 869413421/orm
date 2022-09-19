package dialect

import "reflect"

var dialectsMaps = map[string]Dialect{}

type Dialect interface {
	DataTypeOf(t reflect.Value) string                      // 映射的数据库类型
	TableExistSQL(tableName string) (string, []interface{}) // 返回某个表是否存在的sql语句
}

// RegisterDialect 注册Dialect
func RegisterDialect(name string, dialect Dialect) {
	dialectsMaps[name] = dialect
}

// GetDialect 获取Dialect
func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMaps[name]
	return
}
