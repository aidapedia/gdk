package file

type KeyItf interface {
	Node
	GetValue() interface{}
}

type Key struct {
	name  string
	value interface{}
}

func NewKey(name string, value interface{}) KeyItf {
	return &Key{
		name:  name,
		value: value,
	}
}

func (k *Key) GetValue() interface{} {
	return k.value
}

func (k *Key) GetType() nodeType {
	return nodeTypeKey
}

func (k *Key) GetName() string {
	return k.name
}
