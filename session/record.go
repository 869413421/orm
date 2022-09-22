package session

import (
	"errors"
	"github.com/869413421/orm/clause"
	"reflect"
)

// Insert 插入
func (s *Session) Insert(values ...interface{}) (int64, error) {
	// 1.生成insert语句已经返回所有值
	recordValues := make([]interface{}, 0)
	for _, value := range values {
		s.CallMethod(AfterInsert, value)
		table := s.Model(value).RefTable()
		s.clause.Set(clause.INSERT, table.TableName, table.FieldNames)
		recordValues = append(recordValues, table.RecordValues(value))
	}

	// 2.生成value以及对应值
	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// Find 查找
func (s *Session) Find(values interface{}) error {
	// 1.获取到反射对象
	destSlice := reflect.Indirect(reflect.ValueOf(values))

	// 2.确认单个元素的类型
	destType := destSlice.Type().Elem()

	// 3.创建一个确认的类型拿到类型的table信息
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()

	// 4.生成select 语句并且执行，取到符合条件的结果集
	s.clause.Set(clause.SELECT, table.TableName, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDER, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}

	// 5.遍历结果集
	for rows.Next() {
		// 5.1 创建类型的实例
		dest := reflect.New(destType).Elem()

		// 5.2 将实例的指针传递到values中
		var values []interface{}
		for _, name := range table.FieldNames {
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}

		// 5.3 该行记录每一列的值依次赋值给 values 中的每一个字段
		if err = rows.Scan(values...); err != nil {
			return err
		}

		s.CallMethod(AfterQuery, dest.Addr().Interface())
		// 5.4 将实例添加到切片中
		destSlice.Set(reflect.Append(destSlice, dest))
	}

	return rows.Close()
}

// Update support map[string]interface{}
// also support kv list: "Name", "Tom", "Age", 18, ....
func (s *Session) Update(kv ...interface{}) (int64, error) {
	// 1.检查入参是否为字典，如果不是转换为字典
	m, ok := kv[0].(map[string]interface{})
	if !ok {
		m = make(map[string]interface{})
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
	}
	s.clause.Set(clause.UPDATE, s.refTable.TableName, m)
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// Delete  删除方法
func (s *Session) Delete() (int64, error) {
	s.clause.Set(clause.DELETE, s.refTable.TableName)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// Count  统计count
func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.refTable.TableName)
	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, vars...).QueryRow()
	var tmp int64

	err := row.Scan(&tmp)
	if err != nil {
		return 0, err
	}

	return tmp, err
}

// Limit limit
func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s
}

// Where Where
func (s *Session) Where(desc string, args ...interface{}) *Session {
	var vars []interface{}
	s.clause.Set(clause.WHERE, append(append(vars, desc), args...)...)
	return s
}

// OrderBy adds order by condition to clause
func (s *Session) OrderBy(desc string) *Session {
	s.clause.Set(clause.ORDER, desc)
	return s
}

// First get first
func (s *Session) First(value interface{}) error {
	dest := reflect.Indirect(reflect.ValueOf(value))
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	err := s.Limit(1).Find(destSlice.Addr().Interface())
	if err != nil {
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("NOT FOUND")
	}
	dest.Set(destSlice.Index(0))
	return nil
}
