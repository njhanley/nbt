package nbt

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
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

var typeIDs = map[string]Type{
	"End":       TypeEnd,
	"Byte":      TypeByte,
	"Short":     TypeShort,
	"Int":       TypeInt,
	"Long":      TypeLong,
	"Float":     TypeFloat,
	"Double":    TypeDouble,
	"ByteArray": TypeByteArray,
	"String":    TypeString,
	"List":      TypeList,
	"Compound":  TypeCompound,
	"IntArray":  TypeIntArray,
	"LongArray": TypeLongArray,
}

func (typ Type) String() string {
	if int(typ) >= len(typeNames) {
		return fmt.Sprintf("%#02x", byte(typ))
	}
	return typeNames[typ]
}

func (typ Type) MarshalJSON() ([]byte, error) {
	return json.Marshal(typ.String())
}

func (typ *Type) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	_typ, ok := typeIDs[s]
	if !ok {
		return fmt.Errorf("unknown type (%v)", typ)
	}

	*typ = _typ

	return nil
}

type NamedTag struct {
	Type    Type
	Name    string
	Payload interface{}
}

type jsonNamedTag struct {
	Type    Type            `json:"type"`
	Name    string          `json:"name"`
	Payload json.RawMessage `json:"payload"`
}

func (tag *NamedTag) MarshalJSON() ([]byte, error) {
	payload, err := payloadMarshalJSON(tag.Type, tag.Payload)
	if err != nil {
		return nil, err
	}
	return json.Marshal(&jsonNamedTag{tag.Type, tag.Name, payload})
}

func (tag *NamedTag) UnmarshalJSON(data []byte) error {
	_tag := new(jsonNamedTag)
	if err := json.Unmarshal(data, _tag); err != nil {
		return err
	}

	payload, err := payloadUnmarshalJSON(_tag.Type, _tag.Payload)
	if err != nil {
		return err
	}

	*tag = NamedTag{_tag.Type, _tag.Name, payload}

	return nil
}

func payloadMarshalJSON(typ Type, payload interface{}) (json.RawMessage, error) {
	var _payload interface{}
	switch typ {
	case TypeEnd:
		_payload = payload
	case TypeByte:
		_payload = strconv.FormatInt(int64(payload.(int8)), 10)
	case TypeShort:
		_payload = strconv.FormatInt(int64(payload.(int16)), 10)
	case TypeInt:
		_payload = strconv.FormatInt(int64(payload.(int32)), 10)
	case TypeLong:
		_payload = strconv.FormatInt(int64(payload.(int64)), 10)
	case TypeFloat:
		_payload = strconv.FormatFloat(float64(payload.(float32)), 'g', -1, 32)
	case TypeDouble:
		_payload = strconv.FormatFloat(float64(payload.(float64)), 'g', -1, 64)
	case TypeByteArray:
		_payload = byteArray(payload.([]byte))
	case TypeString:
		_payload = payload.(string)
	case TypeList:
		_payload = payload.(*List)
	case TypeCompound:
		_payload = payload.(Compound)
	case TypeIntArray:
		_payload = intArray(payload.([]int32))
	case TypeLongArray:
		_payload = longArray(payload.([]int64))
	default:
		return nil, fmt.Errorf("unknown type (%v)", typ)
	}

	data, err := json.Marshal(_payload)

	return json.RawMessage(data), err
}

func payloadUnmarshalJSON(typ Type, data json.RawMessage) (interface{}, error) {
	payload, err := interface{}(nil), error(nil)
	switch typ {
	case TypeEnd:
		err = json.Unmarshal(data, &payload)
	case TypeByte, TypeShort, TypeInt, TypeLong, TypeFloat, TypeDouble:
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return nil, err
		}

		n, x := int64(0), float64(0)
		switch typ {
		case TypeByte:
			n, err = strconv.ParseInt(s, 10, 8)
			payload = int8(n)
		case TypeShort:
			n, err = strconv.ParseInt(s, 10, 16)
			payload = int16(n)
		case TypeInt:
			n, err = strconv.ParseInt(s, 10, 32)
			payload = int32(n)
		case TypeLong:
			n, err = strconv.ParseInt(s, 10, 64)
			payload = int64(n)
		case TypeFloat:
			x, err = strconv.ParseFloat(s, 32)
			payload = float32(x)
		case TypeDouble:
			x, err = strconv.ParseFloat(s, 64)
			payload = float64(x)
		}
	case TypeByteArray:
		var b byteArray
		err = json.Unmarshal(data, &b)
		payload = []byte(b)
	case TypeString:
		var s string
		err = json.Unmarshal(data, &s)
		payload = s
	case TypeList:
		l := new(List)
		err = json.Unmarshal(data, l)
		payload = l
	case TypeCompound:
		var m Compound
		err = json.Unmarshal(data, &m)
		payload = m
	case TypeIntArray:
		var a intArray
		err = json.Unmarshal(data, &a)
		payload = []int32(a)
	case TypeLongArray:
		var a longArray
		err = json.Unmarshal(data, &a)
		payload = []int64(a)
	default:
		return nil, fmt.Errorf("unknown type (%v)", typ)
	}

	return payload, err
}

type byteArray []byte

func (b byteArray) MarshalJSON() ([]byte, error) {
	ss := make([]string, len(b))
	for i, n := range b {
		ss[i] = strconv.FormatUint(uint64(n), 10)
	}
	return json.Marshal(ss)
}

func (b *byteArray) UnmarshalJSON(data []byte) error {
	var ss []string
	if err := json.Unmarshal(data, &ss); err != nil {
		return err
	}

	_b := make([]byte, len(ss))
	for i, s := range ss {
		n, err := strconv.ParseUint(s, 10, 8)
		if err != nil {
			return err
		}
		_b[i] = byte(n)
	}

	*b = _b

	return nil
}

type intArray []int32

func (a intArray) MarshalJSON() ([]byte, error) {
	ss := make([]string, len(a))
	for i, n := range a {
		ss[i] = strconv.FormatInt(int64(n), 10)
	}
	return json.Marshal(ss)
}

func (a *intArray) UnmarshalJSON(data []byte) error {
	var ss []string
	if err := json.Unmarshal(data, &ss); err != nil {
		return err
	}

	_a := make([]int32, len(ss))
	for i, s := range ss {
		n, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return err
		}
		_a[i] = int32(n)
	}

	*a = _a

	return nil
}

type longArray []int64

func (a longArray) MarshalJSON() ([]byte, error) {
	ss := make([]string, len(a))
	for i, n := range a {
		ss[i] = strconv.FormatInt(int64(n), 10)
	}
	return json.Marshal(ss)
}

func (a *longArray) UnmarshalJSON(data []byte) error {
	var ss []string
	if err := json.Unmarshal(data, &ss); err != nil {
		return err
	}

	_a := make([]int64, len(ss))
	for i, s := range ss {
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		_a[i] = int64(n)
	}

	*a = _a

	return nil
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

type jsonList struct {
	Type  Type            `json:"type"`
	Array json.RawMessage `json:"array"`
}

func (l *List) MarshalJSON() ([]byte, error) {
	var _array interface{}
	switch l.Type {
	case TypeEnd:
		_array = l.Array
	case TypeByte:
		ss := make([]string, l.Length())
		for i, n := range l.Array.([]int8) {
			ss[i] = strconv.FormatInt(int64(n), 10)
		}
		_array = ss
	case TypeShort:
		ss := make([]string, l.Length())
		for i, n := range l.Array.([]int16) {
			ss[i] = strconv.FormatInt(int64(n), 10)
		}
		_array = ss
	case TypeInt:
		ss := make([]string, l.Length())
		for i, n := range l.Array.([]int32) {
			ss[i] = strconv.FormatInt(int64(n), 10)
		}
		_array = ss
	case TypeLong:
		ss := make([]string, l.Length())
		for i, n := range l.Array.([]int64) {
			ss[i] = strconv.FormatInt(int64(n), 10)
		}
		_array = ss
	case TypeFloat:
		ss := make([]string, l.Length())
		for i, x := range l.Array.([]float32) {
			ss[i] = strconv.FormatFloat(float64(x), 'g', -1, 32)
		}
		_array = ss
	case TypeDouble:
		ss := make([]string, l.Length())
		for i, x := range l.Array.([]float64) {
			ss[i] = strconv.FormatFloat(float64(x), 'g', -1, 64)
		}
		_array = ss
	case TypeByteArray:
		bs := make([]byteArray, l.Length())
		for i, b := range l.Array.([][]byte) {
			bs[i] = byteArray(b)
		}
		_array = bs
	case TypeString:
		_array = l.Array.([]string)
	case TypeList:
		_array = l.Array.([]*List)
	case TypeCompound:
		_array = l.Array.([]Compound)
	case TypeIntArray:
		as := make([]intArray, l.Length())
		for i, a := range l.Array.([][]int32) {
			as[i] = intArray(a)
		}
		_array = as
	case TypeLongArray:
		as := make([]longArray, l.Length())
		for i, a := range l.Array.([][]int64) {
			as[i] = longArray(a)
		}
		_array = as
	default:
		return nil, fmt.Errorf("unknown type (%v)", l.Type)
	}

	array, err := json.Marshal(_array)
	if err != nil {
		return nil, err
	}

	return json.Marshal(&jsonList{l.Type, json.RawMessage(array)})
}

func (l *List) UnmarshalJSON(data []byte) error {
	_l := new(jsonList)
	if err := json.Unmarshal(data, _l); err != nil {
		return err
	}

	var array interface{}
	switch _l.Type {
	case TypeEnd:
		if err := json.Unmarshal(_l.Array, &array); err != nil {
			return err
		}
	case TypeByte, TypeShort, TypeInt, TypeLong, TypeFloat, TypeDouble:
		var ss []string
		if err := json.Unmarshal(_l.Array, &ss); err != nil {
			return err
		}

		switch _l.Type {
		case TypeByte:
			_array := make([]int8, len(ss))
			for i, s := range ss {
				n, err := strconv.ParseInt(s, 10, 8)
				if err != nil {
					return err
				}
				_array[i] = int8(n)
			}
			array = _array
		case TypeShort:
			_array := make([]int16, len(ss))
			for i, s := range ss {
				n, err := strconv.ParseInt(s, 10, 16)
				if err != nil {
					return err
				}
				_array[i] = int16(n)
			}
			array = _array
		case TypeInt:
			_array := make([]int32, len(ss))
			for i, s := range ss {
				n, err := strconv.ParseInt(s, 10, 32)
				if err != nil {
					return err
				}
				_array[i] = int32(n)
			}
			array = _array
		case TypeLong:
			_array := make([]int64, len(ss))
			for i, s := range ss {
				n, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					return err
				}
				_array[i] = int64(n)
			}
			array = _array
		case TypeFloat:
			_array := make([]float32, len(ss))
			for i, s := range ss {
				x, err := strconv.ParseFloat(s, 32)
				if err != nil {
					return err
				}
				_array[i] = float32(x)
			}
			array = _array
		case TypeDouble:
			_array := make([]float64, len(ss))
			for i, s := range ss {
				x, err := strconv.ParseFloat(s, 64)
				if err != nil {
					return err
				}
				_array[i] = float64(x)
			}
			array = _array
		}
	case TypeByteArray:
		var bs []byteArray
		if err := json.Unmarshal(_l.Array, &bs); err != nil {
			return err
		}

		_array := make([][]byte, len(bs))
		for i, b := range bs {
			_array[i] = []byte(b)
		}

		array = _array
	case TypeString:
		var _array []string
		if err := json.Unmarshal(_l.Array, &_array); err != nil {
			return err
		}
		array = _array
	case TypeList:
		var _array []*List
		if err := json.Unmarshal(_l.Array, &_array); err != nil {
			return err
		}
		array = _array
	case TypeCompound:
		var _array []Compound
		if err := json.Unmarshal(_l.Array, &_array); err != nil {
			return err
		}
		array = _array
	case TypeIntArray:
		var as []intArray
		if err := json.Unmarshal(_l.Array, &as); err != nil {
			return err
		}

		_array := make([][]int32, len(as))
		for i, a := range as {
			_array[i] = []int32(a)
		}

		array = _array
	case TypeLongArray:
		var as []longArray
		if err := json.Unmarshal(_l.Array, &as); err != nil {
			return err
		}

		_array := make([][]int64, len(as))
		for i, a := range as {
			_array[i] = []int64(a)
		}

		array = _array
	default:
		return fmt.Errorf("unknown type (%v)", _l.Type)
	}

	*l = List{_l.Type, array}

	return nil
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

type Compound map[string]*Tag

type Tag struct {
	Type    Type
	Payload interface{}
}

type jsonTag struct {
	Type    Type            `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func (tag *Tag) MarshalJSON() ([]byte, error) {
	payload, err := payloadMarshalJSON(tag.Type, tag.Payload)
	if err != nil {
		return nil, err
	}
	return json.Marshal(&jsonTag{tag.Type, payload})
}

func (tag *Tag) UnmarshalJSON(data []byte) error {
	_tag := new(jsonTag)
	if err := json.Unmarshal(data, _tag); err != nil {
		return err
	}

	payload, err := payloadUnmarshalJSON(_tag.Type, _tag.Payload)
	if err != nil {
		return err
	}

	*tag = Tag{_tag.Type, payload}

	return nil
}
