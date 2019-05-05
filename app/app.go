package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type App struct {
	Tree     []*TreeNode
	TreeMap  map[int]*TreeNode
	TreePath string
}

func initRouters(app *App, router *gin.Engine) {
	router.GET("/tree/get_node", func(c *gin.Context) {
		id := c.Query("id")
		if id == "" {
			id = "/"
		}

		el := app.findTree(id)
		c.JSON(http.StatusOK, el)
	})
	router.GET("/tree/get_content", func(c *gin.Context) {
		var id int
		if v, e := strconv.Atoi(c.Query("id")); e == nil {
			id = v
		}

		el := app.TreeMap[id]
		c.JSON(http.StatusOK, gin.H{
			"content": el.Content,
		})
	})
	router.GET("/tree/create_node", func(c *gin.Context) {
		var id int
		var position int
		if v, e := strconv.Atoi(c.Query("id")); e == nil {
			id = v
		}

		if v, e := strconv.Atoi(c.Query("position")); e == nil {
			position = v
		}

		text := c.Query("text")

		el := app.createNode(id, position, text)
		app.saveTree()
		app.readTree()

		c.JSON(http.StatusOK, gin.H{
			"id": el.ID,
		})
	})
	router.GET("/tree/rename_node", func(c *gin.Context) {
		var id int
		if v, e := strconv.Atoi(c.Query("id")); e == nil {
			id = v
		}

		text := c.Query("text")

		el := app.TreeMap[id]
		el.Name = text

		app.saveTree()
		app.readTree()

		c.JSON(http.StatusOK, true)
	})
	router.GET("/tree/delete_node", func(c *gin.Context) {
		var id int
		if v, e := strconv.Atoi(c.Query("id")); e == nil {
			id = v
		}

		app.Tree = app.deleteTree(app.Tree, id)

		app.saveTree()
		app.readTree()

		c.JSON(http.StatusOK, true)
	})
	router.GET("/tree/move_node", func(c *gin.Context) {
		var id int
		var parent int
		var position int
		if v, e := strconv.Atoi(c.Query("id")); e == nil {
			id = v
		}

		if v, e := strconv.Atoi(c.Query("parent")); e == nil {
			parent = v
		}

		if v, e := strconv.Atoi(c.Query("position")); e == nil {
			position = v
		}

		el := *app.TreeMap[id]
		app.Tree = app.deleteTree(app.Tree, id)

		if c.Query("parent") == "#" {
			app.Tree = append(app.Tree, &TreeNode{} /* use the zero value of the element type */)
			copy(app.Tree[position+1:], app.Tree[position:])
			newEl := &TreeNode{
				ID:      id,
				Name:    el.Name,
				Content: el.Content,
			}
			app.Tree[position] = newEl
		} else {
			node := app.createNode(parent, position, "")
			node.ID = el.ID
			node.Name = el.Name
			node.Content = el.Content
			node.Children = el.Children
		}

		app.saveTree()
		app.readTree()
		c.JSON(http.StatusOK, true)
	})

	router.POST("/save-page", func(c *gin.Context) {
		var id int
		if v, e := strconv.Atoi(c.PostForm("id")); e == nil {
			id = v
		}
		title := c.PostForm("title")
		text := c.PostForm("text")

		el := app.TreeMap[id]
		el.Name = text
		el.Name = title
		el.Content = text

		app.saveTree()
		app.readTree()

		c.JSON(http.StatusOK, true)
	})
}

func (a *App) updateIDForChildren(t []*TreeNode) {
	for i := range t {
		t[i].ID = a.emptyID()
		if t[i].Children == nil {
			continue
		}

		a.updateIDForChildren(t[i].Children)
	}
}

func (a *App) readTree() {
	if _, err := os.Stat(a.TreePath); os.IsNotExist(err) {
		empty := []byte("[]")
		err = ioutil.WriteFile(a.TreePath, empty, 0644)
		if err != nil {
			logrus.Fatal(err)
		}
	}

	file, err := ioutil.ReadFile(a.TreePath)
	if err != nil {
		logrus.Fatal(err)
	}

	err = json.Unmarshal([]byte(file), &a.Tree)
	if err != nil {
		logrus.Fatal(err)
	}

	a.TreeMap = make(map[int]*TreeNode, 0)

	a.readTreeToMap(a.Tree)
}

func (a *App) readTreeToMap(t []*TreeNode) {
	for _, el := range t {
		a.TreeMap[el.ID] = el

		if el.Children == nil {
			continue
		}
		a.readTreeToMap(el.Children)
	}
}

func (a *App) deleteTree(t []*TreeNode, id int) []*TreeNode {
	list := make([]*TreeNode, 0)

	for _, el := range t {
		if el.ID == id {
			continue
		}

		el.Children = a.deleteTree(el.Children, id)
		list = append(list, el)
	}

	return list
}

func (a *App) saveTree() {
	file, err := json.MarshalIndent(a.Tree, "", " ")
	if err != nil {
		logrus.Fatal(err)
	}

	err = ioutil.WriteFile(a.TreePath, file, 0644)
	if err != nil {
		logrus.Fatal(err)
	}
}

func (a *App) findTree(idStr string) []TreeNodePublic {
	list := make([]TreeNodePublic, 0, 0)

	if idStr == "#" {
		for _, t := range a.Tree {
			list = append(list, t.Get())
		}
	} else {
		var id int
		if v, e := strconv.Atoi(idStr); e == nil {
			id = v
		}

		for _, t := range a.TreeMap[id].Children {
			list = append(list, t.Get())
		}
	}

	return list
}

func (a *App) emptyID() int {
	for i := 1; ; i++ {
		if _, ok := a.TreeMap[i]; !ok {
			return i
		}
	}
}

func (a *App) createNode(parendID, position int, name string) *TreeNode {
	id := a.emptyID()

	el := a.TreeMap[parendID]

	if el.Children == nil {
		el.Children = make([]*TreeNode, 0)
	}

	el.Children = append(el.Children, &TreeNode{} /* use the zero value of the element type */)
	copy(el.Children[position+1:], el.Children[position:])
	newEl := &TreeNode{
		ID:   id,
		Name: name,
	}
	el.Children[position] = newEl

	return newEl
}

// Run it is main of programm
func Run(router *gin.Engine) error {
	app := &App{
		TreePath: "./db.json",
	}

	app.readTree()

	// app.Tree = []TreeNode{
	// 	TreeNode{
	// 		ID:   1,
	// 		Name: "1",
	// 		Children: []TreeNode{
	// 			TreeNode{
	// 				ID:   2,
	// 				Name: "1_1",
	// 			},
	// 			TreeNode{
	// 				ID:   3,
	// 				Name: "1_2",
	// 			},
	// 		},
	// 	},
	// 	TreeNode{
	// 		ID:   4,
	// 		Name: "2",
	// 	},
	// 	TreeNode{
	// 		ID:   5,
	// 		Name: "3",
	// 	},
	// 	TreeNode{
	// 		ID:   6,
	// 		Name: "4",
	// 	},
	// }

	// app.saveTree()

	initRouters(app, router)
	return nil
}
