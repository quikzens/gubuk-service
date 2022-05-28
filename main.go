package main

import (
	"embed"
	"gubuk-service/config"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

//go:embed build/*
var client embed.FS

type embedFileSystem struct {
	http.FileSystem
	indexes bool
}

func (e embedFileSystem) Exists(prefix string, path string) bool {
	f, err := e.Open(path)
	if err != nil {
		return false
	}

	// check if indexing is allowed
	s, _ := f.Stat()
	if s.IsDir() && !e.indexes {
		return false
	}

	return true
}

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-type"},
		AllowAllOrigins:  true,
		AllowCredentials: true,
	}))

	// Client
	subFS, _ := fs.Sub(client, "build")
	staticServer := static.Serve("/", embedFileSystem{
		FileSystem: http.FS(subFS),
		indexes:    true,
	})
	router.Use(staticServer)
	router.NoRoute(func(c *gin.Context) {
		if c.Request.Method == http.MethodGet &&
			!strings.ContainsRune(c.Request.URL.Path, '.') &&
			!strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.Request.URL.Path = "/"
			staticServer(c)
		}
	})

	SetRoutes(router)

	log.Fatal(router.Run(config.ServerAddress))
}
