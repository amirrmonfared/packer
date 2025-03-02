package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	_ "github.com/amirrmonfared/packer/docs/packer"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"github.com/amirrmonfared/packer/pkg/server"
	"github.com/amirrmonfared/packer/pkg/store"
)

// @title         Packer Service
// @version       1.0
// @description   A simple pack calculator service.
//
// @BasePath /api/v1

type Options struct {
	port int
}

func gatherOptions() Options {
	o := Options{}
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	fs.IntVar(&o.port, "port", 8080, "Port number where the server will listen")

	if err := fs.Parse(os.Args[1:]); err != nil {
		logrus.WithError(err).Fatal("couldn't parse arguments.")
	}
	return o
}

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	logrus.SetLevel(logrus.InfoLevel)

	opts := gatherOptions()

	packStore := store.NewMemoryStore([]int{250, 500, 1000, 2000, 5000})

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	server.RegisterRoutes(r, packStore)

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Static("/web", "./web")
	r.GET("/", func(c *gin.Context) {
		c.File("./web/index.html")
	})

	logrus.Infof("Starting Pack Calculator on port %d", opts.port)
	if err := r.Run(fmt.Sprintf(":%d", opts.port)); err != nil {
		logrus.WithError(err).Fatal("error occurred while running the api server")
	}
}
