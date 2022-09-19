package main

import (
	"fmt"
	"github.com/869413421/orm"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	engine, _ := orm.NewEngine("sqlite3", "gee.db")
	defer engine.Close()
	s := engine.NewSession()

	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	s.Raw("CREATE TABLE User(Name text)").Exec()
	s.Raw("CREATE TABLE User(Name text)").Exec()
	result, _ := s.Raw("INSERT INTO User(`name`) VALUES (?) , (?)", "TOM", "SAM").Exec()
	count, _ := result.RowsAffected()
	fmt.Printf("Exec success, %d addected \n", count)
}
