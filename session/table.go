package session

import (
	"fmt"
	"geeorm/log"
	"geeorm/schema"
	"reflect"
	"strings"
)

//解析对象为表结构，并存放到session中，避免相同的对象重复解析
func (s *Session) Model(value interface{}) *Session {
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("Model is not set")
	}
	return s.refTable
}

func (s *Session) CreateTable() error {
	table := s.RefTable()
	columns := []string{}
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	desc := strings.Join(columns, ",")
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s);", table.Name, desc)).Exec()
	return err
}

func (s *Session) DropTable() error {
	table := s.RefTable()
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s;", table.Name)).Exec()
	return err
}

func (s *Session) HashTable() bool {
	// table := s.RefTable()
	// row := s.Raw(fmt.Sprintf("SELECT name FROM sqlite_master WHERE type='table' and name = %s;", table.Name)).QueryRow()
	// var tmp string
	// _ = row.Scan(&tmp)
	// fmt.Println(tmp)
	// return table.Name == tmp
	sql, values := s.dialect.TableExistSQL(s.RefTable().Name)
	row := s.Raw(sql, values...).QueryRow()
	var tmp string
	_ = row.Scan(&tmp)
	return tmp == s.RefTable().Name
}
