package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type App struct {
	Three     []ThreeNode
	ThreePath string
}

func initRouters(app *App, router *gin.Engine) {
	router.GET("/api/top", func(c *gin.Context) {
		c.JSON(http.StatusOK, 1)
	})
	router.GET("/three/get_node", func(c *gin.Context) {
		id := c.Param("id")
		el := app.findTree(id)
		c.JSON(http.StatusOK, el)
	})
}

func (a *App) readTree() {
	if _, err := os.Stat(a.ThreePath); os.IsNotExist(err) {
		empty := []byte("[]")
		err = ioutil.WriteFile(a.ThreePath, empty, 0644)
		if err != nil {
			logrus.Fatal(err)
		}
	}

	file, err := ioutil.ReadFile(a.ThreePath)
	if err != nil {
		logrus.Fatal(err)
	}

	err = json.Unmarshal([]byte(file), &a.Three)
	if err != nil {
		logrus.Fatal(err)
	}
}

func (a *App) saveTree() {
	file, err := json.MarshalIndent(a.Three, "", " ")
	if err != nil {
		logrus.Fatal(err)
	}

	err = ioutil.WriteFile(a.ThreePath, file, 0644)
	if err != nil {
		logrus.Fatal(err)
	}
}

func (a *App) findTree(id string) []ThreeNodePublic {
	list := make([]ThreeNodePublic, 0, 0)

	if id == "#" {
		for _, t := range a.Three {
			list = append(list, t.Get())
		}
	}

	return list
}

// Run it is main of programm
func Run(router *gin.Engine) error {
	app := &App{
		ThreePath: "./db.json",
	}

	app.readTree()

	app.Three = []ThreeNode{
		ThreeNode{
			Name: "1",
		},
		ThreeNode{
			Name: "2",
		},
		ThreeNode{
			Name: "3",
		},
		ThreeNode{
			Name: "4",
		},
	}

	initRouters(app, router)
	return nil
}
