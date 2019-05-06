package main

import (
	"military/app"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetReportCaller(true)

	r, err := SetupRouter()
	if err != nil {
		logrus.Panicf("Error init: %s\n", err.Error())
	}

	if r.Run(":8080") != nil {
		logrus.Panicf("Error executing: %s\n", err.Error())
	}
}

// SetupRouter return gin router. Made in a separate function for writing tests.
func SetupRouter() (*gin.Engine, error) {
	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("public", false)))

	var err error

	err = app.Run(router.Group("/BKS/"), "./public/db/BKS.json")
	err = app.Run(router.Group("/SB/"), "./public/db/SB.json")
	err = app.Run(router.Group("/VMF/"), "./public/db/VMF.json")
	err = app.Run(router.Group("/BDB/"), "./public/db/BDB.json")
	err = app.Run(router.Group("/RVSM/"), "./public/db/RVSM.json")
	err = app.Run(router.Group("/GUAP/"), "./public/db/GUAP.json")

	return router, err
}
