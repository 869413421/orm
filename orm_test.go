package orm

import (
	"errors"
	"github.com/869413421/orm/session"
	"testing"
)

type User struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}

func TestEngine_Transaction(t *testing.T) {
	t.Run("rollback", func(t *testing.T) {
		transactionRollback(t)
	})
	t.Run("commit", func(t *testing.T) {
		transactionCommit(t)
	})
}

func OpenDb(t *testing.T) *Engine {
	t.Helper()
	e, err := NewEngine("sqlite3", "gee.db")
	if err != nil {
		t.Fatal("fail to connect", err)
	}
	return e
}

func transactionRollback(t *testing.T) {
	engine := OpenDb(t)
	defer engine.Close()

	s := engine.NewSession()
	_ = s.Model(&User{}).DropTable()

	_, err := engine.Transaction(func(session *session.Session) (result interface{}, err error) {
		_ = session.Model(&User{}).CreateTable()
		_, _ = session.Insert(&User{"Tom", 18})
		return nil, errors.New("Error")
	})

	if err == nil || s.HasTable() {
		t.Fatal("failed to rollback")
	}
}


func transactionCommit(t *testing.T) {
	engine := OpenDb(t)
	defer engine.Close()
	s := engine.NewSession()
	_ = s.Model(&User{}).DropTable()
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		_ = s.Model(&User{}).CreateTable()
		_, err = s.Insert(&User{"Tom", 18})
		return
	})
	u := &User{}
	_ = s.First(u)
	if err != nil || u.Name != "Tom" {
		t.Fatal("failed to commit")
	}
}