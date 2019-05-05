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

	if r.Run(":4200") != nil {
		logrus.Panicf("Error executing: %s\n", err.Error())
	}
}

// SetupRouter return gin router. Made in a separate function for writing tests.
func SetupRouter() (*gin.Engine, error) {
	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("public", false)))

	err := app.Run(router)

	return router, err
}
