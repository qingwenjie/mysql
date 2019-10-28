package mysql

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"xorm.io/core"
)

const (
	MAX_CONNS  = 10
	IDLE_CONNS = 3
)

type Options struct {
	Host        string
	Port        int
	User        string
	Password    string
	Database    string
	Charset     string
	MaxConnect  int
	IdleConnect int
	ShowSql     bool
	TablePrefix string
}

type mysql struct {
	ErrorNew error
}

var db *xorm.Engine

func (s *mysql) connect(options *Options) *xorm.Engine {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s",
		options.User,
		options.Password,
		options.Host,
		options.Port,
		options.Database,
		options.Charset,
	)
	engine, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		s.ErrorNew = errors.New("Mysql Engine can't Create: " + err.Error())
		return nil
	}
	engine.SetMaxOpenConns(options.MaxConnect)
	engine.SetMaxIdleConns(options.IdleConnect)
	engine.ShowSQL(options.ShowSql)
	engine.SetMapper(core.SnakeMapper{})
	engine.SetTableMapper(core.NewPrefixMapper(core.SnakeMapper{}, options.TablePrefix))
	engine.SetLogger(xorm.NewSimpleLogger(&xormLogWriter{}))
	if err = engine.Ping(); err != nil {
		s.ErrorNew = err
		return nil
	}
	return engine
}

//连接数据库
func Connect(options *Options) *xorm.Engine {
	if db != nil {
		return db
	}
	if options == nil {
		options = DefaultConfig()
	}
	if options.Charset == "" {
		options.Charset = "utf8mb4"
	}
	if options.Port == 0 {
		options.Port = 3306
	}
	mysql := mysql{}
	engine := mysql.connect(options)
	if mysql.ErrorNew != nil {
		fmt.Println(mysql.ErrorNew)
		return nil
	} else if engine == nil {
		fmt.Println("mysql connect failed")
		return nil
	}
	db = engine
	return db
}

//重新连接数据库连接
func ReConnect(options *Options) *xorm.Engine {
	if db != nil {
		db.Close()
	}
	return Connect(options)
}

func DefaultConfig() *Options {
	return &Options{
		Host:        "127.0.0.1",
		Port:        3306,
		User:        "root",
		Password:    "",
		Database:    "test",
		Charset:     "utf8mb4",
		MaxConnect:  MAX_CONNS,
		IdleConnect: IDLE_CONNS,
		ShowSql:     false,
		TablePrefix: "tb_",
	}
}

type xormLogWriter struct{}

func (w *xormLogWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
