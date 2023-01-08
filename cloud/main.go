package main

import (
	"cloud/lib"
	"cloud/model/mysql"
	"cloud/router"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	serverConfig := lib.LoadServerConfig()
	mysql.InitDB(serverConfig)
	defer mysql.DB.Close()

	r := router.SetupRoute()

	r.LoadHTMLGlob("view/*")
	r.Static("/static", "./static")


	if err := r.Run(":80"); err != nil {
		log.Fatal("服务器启动失败...")
	}
}
