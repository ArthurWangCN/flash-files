package server

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/ArthurWangCN/flash-files/config"
	"github.com/ArthurWangCN/flash-files/server/controller"
	"github.com/ArthurWangCN/flash-files/server/ws"
	"github.com/gin-gonic/gin"
)

//go:embed frontend/dist/*
var FS embed.FS

func Run() {
	hub := ws.NewHub()
	go hub.Run()

	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	staticFiles, _ := fs.Sub(FS, "frontend/dist")
	router.GET("/ws", func(c *gin.Context) {
		ws.HttpController(c, hub)
	})
	router.GET("/api/v1/addresses", controller.AddressesController)
	router.GET("/api/v1/qrcodes", controller.QrcodesController)
	router.POST("/api/v1/texts", controller.TextController)
	router.POST("/api/v1/files", controller.FilesController)
	router.GET("/uploads/:path", controller.UploadsController)
	router.StaticFS("/static", http.FS(staticFiles))
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/static/") {
			reader, err := staticFiles.Open("index.html")
			if err != nil {
				log.Fatal(err)
			}
			defer reader.Close()
			stat, err := reader.Stat()
			if err != nil {
				log.Fatal(err)
			}
			c.DataFromReader(http.StatusOK, stat.Size(), "text/html;charset=utf-8", reader, nil)
		} else {
			c.Status(http.StatusNotFound)
		}
	})
	router.Run(":" + config.GetPort())
}
