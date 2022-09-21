package session

import (
	"database/sql"
	"fmt"
	"github.com/869413421/orm/clause"
	"github.com/869413421/orm/dialect"
	"github.com/869413421/orm/log"
	"github.com/869413421/orm/schema"
	"reflect"
	"strings"
)

type Session struct {
	db       *sql.DB
	sql      strings.Builder
	sqlVars  []interface{}
	dialect  dialect.Dialect
	refTable *schema.Schema
	clause   clause.Clause
}

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{db: db, dialect: dialect}
}

// Model 设置会话模型
func (s *Session) Model(value interface{}) *Session {
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

// RefTable 返回模型详细信息
func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("Model is not set")
	}
	return s.refTable
}

// CreateTable 表创建
func (s *Session) CreateTable() error {
	// 1.使用表字段拼接成SQL
	tableInfo := s.refTable
	var columns []string
	for _, field := range tableInfo.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}

	desc := strings.Join(columns, ",")
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s);", tableInfo.TableName, desc)).Exec()
	return err
}

// DropTable 表删除
func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", s.refTable.TableName)).Exec()
	return err
}

// HasTable 判断表是否存在
func (s *Session) HasTable() bool {
	sql, value := s.dialect.TableExistSQL(s.refTable.TableName)
	row := s.Raw(sql, value...).QueryRow()
	var temp string
	_ = row.Scan(&temp)
	return temp == s.refTable.TableName
}

// Clear 清空sql字符串和参数
func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
	s.clause = clause.Clause{}
}

// DB 获取sql连接
func (s *Session) DB() *sql.DB {
	return s.db
}

// Raw 写入SQL和参数
func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

// Exec 执行SQL
func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	result, err = s.DB().Exec(s.sql.String(), s.sqlVars...)
	if err != nil {
		log.Error(err)
	}
	return
}

// QueryRow 返回第一行
func (s *Session) QueryRow() (row *sql.Row) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

// QueryRows 返回整个结果集
func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	rows, err = s.DB().Query(s.sql.String(), s.sqlVars...)
	if err != nil {
		log.Error(err)
	}
	return
}
