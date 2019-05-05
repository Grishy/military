package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type App struct{}

func initRouters(app *App, router *gin.Engine) {

	router.GET("/api/top", func(c *gin.Context) {
		type block struct {
			Name    string
			Count   int
			Сaption string
		}

		result := []block{
			{
				Name:    "Rack",
				Count:   1,
				Сaption: "50% Full",
			},
			{
				Name:    "Server",
				Count:   5,
				Сaption: "HP - 100%",
			},
			{
				Name:    "CPU",
				Count:   1 + 1 + 4 + 1 + 2,
				Сaption: "60 Cores and 112 Threads",
			},
			{
				Name:    "GPU",
				Count:   0,
				Сaption: "Not found",
			},
		}
		c.JSON(http.StatusOK, result)
	})
}

// Run it is main of programm
func Run(router *gin.Engine) error {
	app := &App{}

	initRouters(app, router)
	return nil
}
