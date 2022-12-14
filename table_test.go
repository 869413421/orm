package orm

import "testing"


func TestSession_CreateTable(t *testing.T) {
	e, _ := NewEngine("sqlite3", "gee.db")
	s := e.NewSession().Model(&User{})
	_ = s.DropTable()
	_ = s.CreateTable()
	if !s.HasTable() {
		t.Fatal("Failed to create table User")
	}
}
