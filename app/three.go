package app

type TreeNodePublic struct {
	ID       int    `json:"id"`
	Name     string `json:"text"`
	Children bool   `json:"children"`
}

type TreeNode struct {
	ID       int
	Name     string
	Content  string
	Children []*TreeNode
}

func (t *TreeNode) Get() TreeNodePublic {
	children := false
	if len(t.Children) > 0 {
		children = true
	}

	return TreeNodePublic{
		ID:       t.ID,
		Name:     t.Name,
		Children: children,
	}
}
