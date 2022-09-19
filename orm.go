package orm

import (
	"database/sql"
	"github.com/869413421/orm/dialect"
	"github.com/869413421/orm/log"
	"github.com/869413421/orm/session"
	_ "github.com/mattn/go-sqlite3"
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
