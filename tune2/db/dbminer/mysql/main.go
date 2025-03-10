package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"tune/db/dbminer"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// type DBMiner interface {
// 	getSchema() //fix type soon
// }

type SQLMiner struct {
	Host   string
	Client *sql.DB
}

func (sqlMiner *SQLMiner) GetSchema() (*dbminer.Schema, error) {
	schema := dbminer.Schema{Name: sqlMiner.Host}
	qDB, err := sqlMiner.Client.Query("SHOW DATABASES where `Database` not in ('sys', 'mysql', 'information_schema', 'performance_schema')")
	if err != nil {
		log.Panic(err)
	}
	defer qDB.Close()

	databaseNames := []string{}
	// GET DB NAMES
	for qDB.Next() {
		storer := sql.RawBytes{}
		err := qDB.Scan(&storer)
		if err != nil {
			log.Panic(err)
		}
		databaseNames = append(databaseNames, string(storer))
	}

	for _, dbName := range databaseNames {
		database := dbminer.Database{Name: dbName}
		dsn := fmt.Sprintf("%s%s", schema.Name, dbName)
		fmt.Println(dsn)
		gormDB, _ := gorm.Open(mysql.Open(dsn))
		db, err := gormDB.DB()
		if err != nil {
			log.Panic(err)
		}
		defer db.Close()

		qTables, err := db.Query("SHOW TABLES")
		if err != nil {
			log.Panic(err)
		}
		defer qTables.Close()

		var tableNames []string
		// GET TABLES
		for qTables.Next() {
			tempNames := sql.RawBytes{}
			qTables.Scan(&tempNames)
			fmt.Printf("%s\n", tempNames)
			tableNames = append(tableNames, string(tempNames))
		}
		for _, tableName := range tableNames {
			tables := dbminer.Table{Name: tableName}
			fmt.Println(fmt.Sprintf("select * from %v", tableName))
			sqlCols, err := db.Query(fmt.Sprintf("select * from %v", tableName))
			if err != nil {
				log.Panic(err)
			}
			defer sqlCols.Close()

			// for sqlCols.Next() {
			// 	colRaw := sql.RawBytes{}
			// 	sqlCols.Scan(&colRaw)
			// 	tables.Columns = append(tables.Columns, string(colRaw))
			// }
			cols, err := sqlCols.Columns()
			if err != nil {
				log.Panic(err)
			}

			tables.Columns = cols
			database.Tables = append(database.Tables, tables)
		}
		// for _, tableNames := qTables {
		// }
		// if err != nil {
		// 	log.Panic(err)
		// }
		schema.Databases = append(schema.Databases, database)
	}
	return &schema, nil
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
	gDB, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		log.Panic(err, "err")
	}
	client, err := gDB.DB()
	if err != nil {
		log.Panic(err)
	}
	defer client.Close()

	sqlMiner := SQLMiner{dsn, client}
	dbminer.Search(&sqlMiner)

	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/store", dbInfo.User, dbInfo.Pass, dbInfo.Host, dbInfo.Port)
	// fmt.Println(dsn)
	// db, _ := gorm.Open(mysql.Open(dsn))
	// // ctx
	// sqldb, _ := db.DB()
	// defer sqldb.Close()
	// // db.close?
	// // allDataQ := "SELECT table_schema, table_name, column_name FROM information_schema.columns WHERE table_schema NOT IN ('sys', 'mysql', 'information_schema', 'performance_schema')  ORDER BY table_schema, table_name"
	// // allDatabaseQ := "SHOW DATABASES WHERE `database` NOT IN ('sys', 'mysql', 'information_schema', 'performance_schema') "
	// rows, err := sqldb.Query("select * from transactions")
	// // db.Table("information_schema.columns").
	// // 	Select("table_schema, table_name, column_name").
	// // 	Where("table_schema NOT IN (?)", []string{"mysql", "information_schema", "performance_schema"}).
	// // 	Order("table_schema, table_name").
	// // 	Scan(&raw)
	// if err != nil {
	// 	log.Panic(err)
	// }
	// defer rows.Close()

	// col, err := rows.Columns()
	// if err != nil {
	// 	log.Panic(err)
	// }
	// // var (
	// // 	col1, col2, col3 sql.RawBytes
	// // )
	// // columns, _ := rows.Columns()
	// // col1 := make([]any, len(columns))
	// col1 := []any{}
	// for range col {
	// 	col1 = append(col1, new(sql.RawBytes))
	// }
	// for rows.Next() {
	// 	// rows.Scan(&col1, &col2, &col3)
	// 	// fmt.Printf("%s,%s,%s\n", col, col2, col3)
	// 	err := rows.Scan(col1...)
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}
	// 	// fmt.Printf("%s\n", col1)
	// }

	// for _, colData := range col1 {
	// 	if raw, ok := colData.(*sql.RawBytes); ok {
	// 		fmt.Println(string(*raw))
	// 	}
	// }
	// // for _, colType := range col {
	// // 	fmt.Println(colType)
	// // 	// sqltx, _ := sqldb.Begin()
	// // 	// sqltx.
	// // 	// raw := []sql.RawBytes{}
	// // 	// db.Raw("select ? from information_schema.columns where ? not in ('mysql', 'information_schema', 'performance_schema')", colType, colType).Scan(&raw)
	// // 	// fmt.Printf("%+s", raw)
	// // 	// // fmt.Println()
	// // }

}
