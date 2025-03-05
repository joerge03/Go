package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBMiner interface {
	getSchema() //fix type soon
}

type Schema struct {
	Host   string
	Client gorm.DB
	ctx    context.Context
}

type DBInfo struct {
	User string
	Pass string
	Port int
	Host string
}

var dbInfo DBInfo

func init() {
	flag.StringVar(&dbInfo.User, "user", "root", "Username of your db, Default: 'root' ")
	flag.StringVar(&dbInfo.Pass, "pass", "", "Password of your db if there is one ")
	flag.IntVar(&dbInfo.Port, "p", 3306, "Port of your desired DB, Default: '3306'")
	flag.StringVar(&dbInfo.Host, "h", "localhost", "Host of your DB, Default: localhost")
	flag.Parse()

	if dbInfo.Port < 0 || dbInfo.Port > 65535 {
		log.Panic("Invalid port")
	}
}

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/", dbInfo.User, dbInfo.Pass, dbInfo.Host, dbInfo.Port)

	fmt.Println(dsn)

	db, _ := gorm.Open(mysql.Open(dsn))
	// ctx
	sqldb, _ := db.DB()

	// db.close?
	allDataQuery := "SELECT table_schema, table_name, column_name FROM information_schema.columns WHERE table_schema NOT IN ('mysql', 'information_schema', 'performance_schema') ORDER BY table_schema, table_name"
	allData, err := sqldb.Query(allDataQuery)
	if err != nil {
		log.Panic(err)
	}

	col, err := allData.ColumnTypes()
	if err != nil {
		log.Panic(err)
	}

	for _, colType := range col {
		fmt.Println(colType.ScanType())
	}

}
