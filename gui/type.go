package gui

type JSONType int

const (
	Root JSONType = iota + 1
	Object
	Array
	Key
	Value
)

var jsonTypeMap = map[JSONType]string{
	Object: "object",
	Array:  "array",
	Key:    "key",
	Value:  "value",
}

func (t JSONType) String() string {
	return jsonTypeMap[t]
}

type ValueType int

const (
	Int ValueType = iota + 1
	String
	Float
	Boolean
)

var valueTypeMap = map[ValueType]string{
	Int:     "int",
	String:  "string",
	Float:   "float",
	Boolean: "boolean",
}

func (v ValueType) String() string {
	return valueTypeMap[v]
}

type Reference struct {
	JSONType  JSONType
	ValueType ValueType
}
