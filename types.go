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

func (tag *NamedTag) ToByte() int8 {
	return tag.Payload.(int8)
}

func (tag *NamedTag) ToShort() int16 {
	return tag.Payload.(int16)
}

func (tag *NamedTag) ToInt() int32 {
	return tag.Payload.(int32)
}

func (tag *NamedTag) ToLong() int64 {
	return tag.Payload.(int64)
}

func (tag *NamedTag) ToFloat() float32 {
	return tag.Payload.(float32)
}

func (tag *NamedTag) ToDouble() float64 {
	return tag.Payload.(float64)
}

func (tag *NamedTag) ToByteArray() []byte {
	return tag.Payload.([]byte)
}

func (tag *NamedTag) ToString() string {
	return tag.Payload.(string)
}

func (tag *NamedTag) ToList() *List {
	return tag.Payload.(*List)
}

func (tag *NamedTag) ToCompound() Compound {
	return tag.Payload.(Compound)
}

func (tag *NamedTag) ToIntArray() []int32 {
	return tag.Payload.([]int32)
}

func (tag *NamedTag) ToLongArray() []int64 {
	return tag.Payload.([]int64)
}

type List struct {
	Type  Type
	Array interface{}
}

func (list *List) ToByte() []int8 {
	return list.Array.([]int8)
}

func (list *List) ToShort() []int16 {
	return list.Array.([]int16)
}

func (list *List) ToInt() []int32 {
	return list.Array.([]int32)
}

func (list *List) ToLong() []int64 {
	return list.Array.([]int64)
}

func (list *List) ToFloat() []float32 {
	return list.Array.([]float32)
}

func (list *List) ToDouble() []float64 {
	return list.Array.([]float64)
}

func (list *List) ToByteArray() [][]byte {
	return list.Array.([][]byte)
}

func (list *List) ToString() []string {
	return list.Array.([]string)
}

func (list *List) ToList() []*List {
	return list.Array.([]*List)
}

func (list *List) ToCompound() []Compound {
	return list.Array.([]Compound)
}

func (list *List) ToIntArray() [][]int32 {
	return list.Array.([][]int32)
}

func (list *List) ToLongArray() [][]int64 {
	return list.Array.([][]int64)
}

type Compound map[string]*NamedTag
