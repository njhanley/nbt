package nbt

import (
	"fmt"
	"reflect"
)

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
	if int(typ) >= len(typeNames) {
		return fmt.Sprintf("%#02x", byte(typ))
	}
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

func (l *List) Length() int {
	if l.Type == TypeEnd && l.Array == nil {
		return 0
	}
	return reflect.ValueOf(l.Array).Len()
}

func (l *List) ToByte() []int8 {
	return l.Array.([]int8)
}

func (l *List) ToShort() []int16 {
	return l.Array.([]int16)
}

func (l *List) ToInt() []int32 {
	return l.Array.([]int32)
}

func (l *List) ToLong() []int64 {
	return l.Array.([]int64)
}

func (l *List) ToFloat() []float32 {
	return l.Array.([]float32)
}

func (l *List) ToDouble() []float64 {
	return l.Array.([]float64)
}

func (l *List) ToByteArray() [][]byte {
	return l.Array.([][]byte)
}

func (l *List) ToString() []string {
	return l.Array.([]string)
}

func (l *List) ToList() []*List {
	return l.Array.([]*List)
}

func (l *List) ToCompound() []Compound {
	return l.Array.([]Compound)
}

func (l *List) ToIntArray() [][]int32 {
	return l.Array.([][]int32)
}

func (l *List) ToLongArray() [][]int64 {
	return l.Array.([][]int64)
}

type Compound map[string]*NamedTag
