package session

import "github.com/869413421/orm/log"

// Begin 开启事务
func (s *Session) Begin() (err error) {
	log.Info("transaction begin")
	s.tx, err = s.db.Begin()
	if err != nil {
		log.Error(err)
		return err
	}
	return
}

// Commit 提交事务
func (s *Session) Commit() (err error) {
	log.Info("transaction commit")
	err = s.tx.Commit()
	if err != nil {
		log.Error(err)
		return err
	}
	return
}

// RollBack 回滚事务
func (s *Session) RollBack() (err error) {
	log.Info("transaction rollback")
	err = s.tx.Rollback()
	if err != nil {
		log.Error(err)
		return err
	}
	return
}
