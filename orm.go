package orm

import (
	"database/sql"
	"fmt"
	"github.com/869413421/orm/dialect"
	"github.com/869413421/orm/log"
	"github.com/869413421/orm/session"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

// NewEngine 初始化引擎
func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}

	err = db.Ping()
	if err != nil {
		log.Error(err)
		return
	}

	dial, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s Not Found", driver)
		return
	}
	e = &Engine{db: db, dialect: dial}
	log.Info("Connect database success")
	return
}

// Close 关闭数据库连接
func (e *Engine) Close() {
	err := e.db.Close()
	if err != nil {
		log.Error("Failed to close database")
	}
	log.Info("Close database success")
}

// NewSession 创建新会话
func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dialect)
}

type TxFunc func(*session.Session) (interface{}, error)

// Transaction 执行事务
func (e *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := e.NewSession()
	err = s.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = s.RollBack()
			panic(p)
		} else if err != nil {
			_ = s.RollBack()
		} else {
			err = s.Commit()
		}
	}()

	return f(s)
}

// difference 获取两个切片差集
func difference(a []string, b []string) (diff []string) {
	mapB := make(map[string]bool)
	for _, v := range b {
		mapB[v] = true
	}

	for _, v := range a {
		if _, ok := mapB[v]; !ok {
			diff = append(diff, v)
		}
	}

	return
}

func (e *Engine) Migrate(value interface{}) error {
	_, err := e.Transaction(func(s *session.Session) (result interface{}, err error) {
		// 1.如果表不存在则创建
		if !s.Model(value).HasTable() {
			log.Infof("table %s doesnt exist", s.RefTable().TableName)
		}
		table := s.RefTable()

		rows, _ := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", table.TableName)).QueryRows()
		columns, _ := rows.Columns()
		addCols := difference(table.FieldNames, columns)
		delCols := difference(columns, table.FieldNames)
		log.Infof("added cols %v, deleted cols %v", addCols, delCols)

		for _, col := range addCols {
			f := table.GetField(col)
			sqlStr := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", table.TableName, f.Name, f.Type)
			if _, err = s.Raw(sqlStr).Exec(); err != nil {
				return
			}
		}

		if len(delCols) == 0 {
			return
		}
		tmp := "tmp_" + table.TableName
		fieldStr := strings.Join(table.FieldNames, ", ")
		s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s from %s;", tmp, fieldStr, table.TableName))
		s.Raw(fmt.Sprintf("DROP TABLE %s;", table.TableName))
		s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", tmp, table.TableName))
		_, err = s.Exec()
		return
	})

	return err
}
