package main

import (
	"fmt"
	"github.com/DualVectorFoil/Zelda/conf"
	"github.com/DualVectorFoil/Zelda/db"
)

func main() {
	mysqlInfo := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true", conf.MYSQL_USERNAME, conf.MYSQL_PASSWORD, conf.MYSQL_IP, conf.MYSQL_PORT, conf.MYSQL_DBNAME)
	db.InitDB(mysqlInfo)
}
