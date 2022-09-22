package session

import (
	"fmt"
	"github.com/869413421/orm/log"
	"reflect"
)

const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
)

type IBeforeQuery interface {
	BeforeQuery(s *Session) error
}

type IAfterQuery interface {
	AfterQuery(s *Session) error
}

type IBeforeInsert interface {
	BeforeInsert(s *Session) error
}

// CallMethod 如果实现了结构体实现某个Hooks接口则调用
func (s *Session) CallMethod(method string, value interface{}) {
	if value == nil {
		return
	}
	var err error
	obj := reflect.ValueOf(value)
	switch method {
	case BeforeQuery:
		if i, ok := obj.Interface().(IBeforeQuery); ok {
			err = i.BeforeQuery(s)
		}
	case AfterQuery:
		if i, ok := obj.Interface().(IAfterQuery); ok {
			err = i.AfterQuery(s)
		}
	case AfterInsert:
		if i, ok := obj.Interface().(IBeforeInsert); ok {
			err = i.BeforeInsert(s)
		}
	default:
		panic("interface not found")
	}
	if err != nil {
		log.Error(fmt.Sprintf("CallMethod error %s", err))
	}
}
