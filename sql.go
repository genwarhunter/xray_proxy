package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var mySqlDB *sql.DB

func CreateConnectMysql() {
	var err error
	mySqlDB, err = sql.Open("mysql", AppConfig.DataSourceName)
	if err != nil {
		log.Fatalln("sql.Open", err)
	}
	log.Println("Соединение с бд установлено")
	return
}

func selectFromPackage() ([]infoPackageRow, bool) {
	var ans []infoPackageRow
	var results *sql.Rows
	results, err := mySqlDB.Query("SELECT p.id, p.name, p.link, p.use FROM Packages p where p.load = 1")
	if err != nil {
		log.Println("selectFromPackage", err)
		return []infoPackageRow{}, false
	}
	for results.Next() {
		var q infoPackageRow
		err = results.Scan(&q.Id, &q.Name, &q.Url, &q.Use)
		if err != nil {
			log.Println("selectFromPackage", err)
			return []infoPackageRow{}, false
		}
		ans = append(ans, q)
	}
	_ = results.Close()
	log.Println(ans)
	return ans, true
}
