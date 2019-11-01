package gui

type Type int

const (
	Unknown Type = iota
	Root
	Object
	Array
	Key
	Value
)

var TypeMap = map[Type]string{
	Unknown: "unknown",
	Root:    "root",
	Object:  "object",
	Array:   "array",
	Key:     "key",
	Value:   "value",
}

func (t Type) String() string {
	return TypeMap[t]
}
