package app

import (
	uuid "github.com/satori/go.uuid"
)

type ThreeNodePublic struct {
	ID       string
	Text     string
	Name     string
	Children bool
}

type ThreeNode struct {
	ID       string
	Name     string
	Text     string
	Children []*ThreeNode
}

func New(parent *ThreeNode, name string) *ThreeNode {
	id := uuid.NewV4()

	return &ThreeNode{
		ID:       id.String(),
		Name:     name,
		Text:     "",
		Children: make([]*ThreeNode, 0, 0),
	}
}

func (t *ThreeNode) Get() ThreeNodePublic {
	children := false
	if len(t.Children) > 0 {
		children = true
	}

	return ThreeNodePublic{
		ID:       t.ID,
		Name:     t.Name,
		Text:     t.Text,
		Children: children,
	}
}
