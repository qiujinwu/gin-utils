package utils

import (
	"database/sql"
	"github.com/codegangsta/cli"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"strconv"
	"time"
	"github.com/jinzhu/gorm"
	"log"
)

func getDbConfigFlag(c *cli.Context) bool{
	for _,value := range c.FlagNames() {
		if value == "mysql_user" {
			return true
		}
	}
	return false
}

func getDbConfig(c *cli.Context,driver string)(config string){
	if getDbConfigFlag(c){
		log.Print("init " + driver + " " + c.String("mysql_user") + ":**@tcp(" +
			c.String("mysql_host") + ":" +
			strconv.Itoa(c.Int("mysql_port")) + ")/" +
			c.String("mysql_db"))

		config = c.String("mysql_user") + ":" +
			c.String("mysql_password") + "@tcp(" +
			c.String("mysql_host") + ":" +
			strconv.Itoa(c.Int("mysql_port")) + ")/" +
			c.String("mysql_db") + "?charset=utf8&parseTime=true"
	}else{
		log.Print("init " + driver + " " +  c.GlobalString("mysql_user") + ":**@tcp(" +
			c.GlobalString("mysql_host") + ":" +
			strconv.Itoa(c.GlobalInt("mysql_port")) + ")/" +
			c.GlobalString("mysql_db"))

		config = c.GlobalString("mysql_user") + ":" +
			c.GlobalString("mysql_password") + "@tcp(" +
			c.GlobalString("mysql_host") + ":" +
			strconv.Itoa(c.GlobalInt("mysql_port")) + ")/" +
			c.GlobalString("mysql_db") + "?charset=utf8&parseTime=true"
	}
	return
}

func OpenDB(c *cli.Context) (db *sql.DB,driver string) {
	driver = "mysql"
	var config string = getDbConfig(c,driver)
	db, err := sql.Open(driver, config)
	if err != nil {
		log.Fatal(err)
		log.Fatalln("连接数据库失败")
	}
	if driver == "mysql" {
		// per issue https://github.com/go-sql-driver/mysql/issues/257
		db.SetMaxIdleConns(0)
	}

	if err := pingDatabase(db); err != nil {
		log.Fatal(err)
		log.Fatalln("ping 数据库" + driver + "失败")
	}
	return
}


func OpenGorm(c *cli.Context,driver string)(db *gorm.DB){
	var config string = getDbConfig(c,driver)
	db, err := gorm.Open("mysql", config)
	if err != nil {
		log.Fatal(err)
	}
	if c.GlobalBool("debug"){
		db.LogMode(true)
	}
	return
}

func pingDatabase(db *sql.DB) (err error) {
	for i := 0; i < 10; i++ {
		err = db.Ping()
		if err == nil {
			return
		}
		log.Print("ping 数据库失败, 1s后重试")
		time.Sleep(time.Second)
	}
	return
}
