package template

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"time"
	"xorm.io/xorm"
)

type Mysql struct {
	User     string `toml:"user"`
	Password string `toml:"password"`
	Host     string `toml:"host"`
	Database string `toml:"database"`
	Charset  string `toml:"charset"`
}

func (msql *Mysql) Engine() (*xorm.Engine, error) {
	Log.Info("Mysql Connection : "+fmt.Sprintf("%+v", msql), zap.String("middleware", "Mysql"))
	dbUrl := fmt.Sprintf("%v:%v@(%v)/%v?charset=%v", msql.User, msql.Password, msql.Host, msql.Database, msql.Charset)
	dbEngine, err := xorm.NewEngine("mysql", dbUrl)
	if err != nil {
		return nil, err
	}
	dbEngine.DatabaseTZ = time.Local
	dbEngine.TZLocation = time.Local
	dbEngine.ShowSQL(true)
	dbEngine.SetConnMaxLifetime(60 * time.Second)
	dbEngine.SetMaxOpenConns(50)
	return dbEngine, err
}
