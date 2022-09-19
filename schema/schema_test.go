package schema

import (
	"github.com/869413421/orm/dialect"
	"testing"
)

type User struct {
	Name string `orm:"PRIMARY KEY"`
	Age  int
}

func TestParse(t *testing.T) {
	user := &User{}
	testDialect, _ := dialect.GetDialect("sqlite3")
	schema := Parse(user, testDialect)
	if schema.TableName != "User" || len(schema.Fields) != 2 {
		t.Fatal("failed to parse User struct")
	}
	if schema.GetField("Name").Tag != "PRIMARY KEY" {
		t.Fatal("failed to parse primary key")
	}
}


