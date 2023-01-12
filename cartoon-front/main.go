package main

import "net/http"

// "net/http"

// "github.com/gin-gonic/gin"

func main() {
	// r := gin.Default()
	// r.LoadHTMLGlob("../view/*.html")
	// r.Static("/static", "../static")
	// r.GET("/", func(ctx *gin.Context) {
	// 	ctx.HTML(http.StatusOK, "index.html", nil)
	// })
	// r.GET("/detail", func(ctx *gin.Context) {
	// 	ctx.HTML(http.StatusOK, "detail.html", nil)
	// })
	// r.Run(":9090")
	http.ListenAndServe(":9090", http.FileServer(http.Dir("../")))
}