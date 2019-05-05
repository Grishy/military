package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	scribble "github.com/nanobox-io/golang-scribble"
)

type App struct {
	Three ThreeNode
}

// a fish
type Fish struct{ Name string }

func initRouters(app *App, router *gin.Engine) {
	router.GET("/api/top", func(c *gin.Context) {
		c.JSON(http.StatusOK, 1)
	})
	router.GET("/three/get_node", func(c *gin.Context) {
		c.JSON(http.StatusOK, 1)
	})
}

func initDB() {
	dir := "./db"

	db, err := scribble.New(dir, nil)
	if err != nil {
		fmt.Println("Error", err)
	}

	// Write a fish to the database
	for _, name := range []string{"onefish", "twofish", "redfish", "bluefish"} {
		db.Write("fish", name, Fish{Name: name})
	}

	// Read a fish from the database (passing fish by reference)
	onefish := Fish{}
	if err := db.Read("fish", "onefish", &onefish); err != nil {
		fmt.Println("Error", err)
	}

	// Read all fish from the database, unmarshaling the response.
	records, err := db.ReadAll("fish")
	if err != nil {
		fmt.Println("Error", err)
	}

	fishies := []Fish{}
	for _, f := range records {
		fishFound := Fish{}
		if err := json.Unmarshal([]byte(f), &fishFound); err != nil {
			fmt.Println("Error", err)
		}
		fishies = append(fishies, fishFound)
	}
}

// Run it is main of programm
func Run(router *gin.Engine) error {
	app := &App{}

	initDB()
	initRouters(app, router)
	return nil
}
