package template

import (
	"testing"
)

func TestMysql(t *testing.T) {
	tem := &Mysql{
		User:     "root",
		Password: "123456",
		Host:     "10.0.3.191:8066",
		Database: "maoti_app_abroad",
		Charset:  "utf8mb4",
	}
	db, err := tem.Engine()
	if err != nil {
		t.Fatal(err)
	}
	db.ShowSQL(true)
	err = db.Ping()
	if err != nil {
		t.Fatal(err)
	}
}
