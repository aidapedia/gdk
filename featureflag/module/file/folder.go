package file

type FolderItf interface {
	Node
	GetChild(name string) Node
	Add(child Node) FolderItf
}

type Folder struct {
	name     string
	children map[string]Node
}

func NewFolder(name string) FolderItf {
	return &Folder{
		name:     name,
		children: map[string]Node{},
	}
}

func (f *Folder) Add(child Node) FolderItf {
	f.children[child.GetName()] = child
	return f
}

func (f *Folder) GetName() string {
	return f.name
}

func (f *Folder) GetType() nodeType {
	return nodeTypeFolder
}

func (f *Folder) GetChild(name string) Node {
	return f.children[name]
}
