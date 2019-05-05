package app

type ThreeNodePublic struct {
	ID       string
	Content  string
	Name     string
	Children bool
}

type ThreeNode struct {
	Name     string
	Content  string
	Parent   *ThreeNode
	Children []*ThreeNode
}

func New(name string, content string) *ThreeNode {
	return &ThreeNode{
		Name:     name,
		Content:  "",
		Children: make([]*ThreeNode, 0, 0),
	}
}

func (t *ThreeNode) Get() ThreeNodePublic {
	children := false
	if len(t.Children) > 0 {
		children = true
	}

	return ThreeNodePublic{
		ID:       t.GetID(),
		Name:     t.Name,
		Children: children,
	}
}

func (t *ThreeNode) GetID() string {
	if t.Parent == nil {
		return "/"
	}

	return t.Parent.GetID() + "/" + t.GetID()
}
