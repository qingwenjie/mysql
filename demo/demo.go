package main

import (
	"fmt"
	"github.com/qingwenjie/mysql"
)

func main() {
	mysql.Connect(&mysql.Options{
		Host:        "127.0.0.1",
		Port:        3306,
		User:        "root",
		Password:    "123456",
		Database:    "task",
		Charset:     "utf8mb4",
		MaxConnect:  10,
		IdleConnect: 3,
		ShowSql:     true,
		TablePrefix: "tb_",
	})
	t()
}

func t() {
	engine := mysql.Connect(nil)
	result, err := engine.Query("show databases")
	fmt.Println(string(result[0]["Database"]), err)
}
