package nbt

type Type byte

const (
	TypeEnd Type = iota
	TypeByte
	TypeShort
	TypeInt
	TypeLong
	TypeFloat
	TypeDouble
	TypeByteArray
	TypeString
	TypeList
	TypeCompound
	TypeIntArray
	TypeLongArray
)

var typeNames = []string{
	TypeEnd:       "End",
	TypeByte:      "Byte",
	TypeShort:     "Short",
	TypeInt:       "Int",
	TypeLong:      "Long",
	TypeFloat:     "Float",
	TypeDouble:    "Double",
	TypeByteArray: "ByteArray",
	TypeString:    "String",
	TypeList:      "List",
	TypeCompound:  "Compound",
	TypeIntArray:  "IntArray",
	TypeLongArray: "LongArray",
}

func (typ Type) String() string {
	return typeNames[typ]
}

type NamedTag struct {
	Type    Type
	Name    string
	Payload interface{}
}

type List struct {
	Type  Type
	Array interface{}
}

type Compound map[string]*NamedTag
