package nbt

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
)

var (
	ErrNegativeLength = errors.New("length is negative")
	ErrUnknownTag     = errors.New("unknown tag ID")
	ErrInvalidRoot    = errors.New("root tag is not compound")
	ErrDuplicateName  = errors.New("duplicate name in compound")
)

type TagType byte

const (
	TypeEnd TagType = iota
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

var tagTypeName = []string{
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

func (typ TagType) String() string {
	return tagTypeName[typ]
}

func (typ TagType) MarshalJSON() ([]byte, error) {
	return json.Marshal(typ.String())
}

type (
	End       struct{}
	Byte      int8
	Short     int16
	Int       int32
	Long      int64
	Float     float32
	Double    float64
	ByteArray []Byte
	String    string
	List      struct {
		ElementType TagType
		Payload     interface{}
	}
	Compound  map[String]NamedTag
	IntArray  []Int
	LongArray []Long
)

func (m Compound) MarshalJSON() ([]byte, error) {
	a := make([]NamedTag, 0, len(m))
	for _, tag := range m {
		a = append(a, tag)
	}
	return json.Marshal(a)
}

func read(r io.Reader, v interface{}) error {
	return binary.Read(r, binary.BigEndian, v)
}

func Decode(r io.Reader) (NamedTag, error) {
	return readNamedTag(r)
}

type NamedTag struct {
	Type    TagType
	Name    String
	Payload interface{}
}

func readNamedTag(r io.Reader) (NamedTag, error) {
	var tag NamedTag

	err := read(r, &tag.Type)
	if err != nil {
		return tag, err
	}

	if tag.Type == TypeEnd {
		return tag, nil
	}

	tag.Name, err = readString(r)
	if err != nil {
		return tag, err
	}

	switch tag.Type {
	case TypeByte:
		var v Byte
		if err := read(r, &v); err != nil {
			return tag, err
		}
		tag.Payload = v
	case TypeShort:
		var v Short
		if err := read(r, &v); err != nil {
			return tag, err
		}
		tag.Payload = v
	case TypeInt:
		var v Int
		if err := read(r, &v); err != nil {
			return tag, err
		}
		tag.Payload = v
	case TypeLong:
		var v Long
		if err := read(r, &v); err != nil {
			return tag, err
		}
		tag.Payload = v
	case TypeFloat:
		var v Float
		if err := read(r, &v); err != nil {
			return tag, err
		}
		tag.Payload = v
	case TypeDouble:
		var v Double
		if err := read(r, &v); err != nil {
			return tag, err
		}
		tag.Payload = v
	case TypeByteArray:
		v, err := readByteArray(r)
		if err != nil {
			return tag, err
		}
		tag.Payload = v
	case TypeString:
		v, err := readString(r)
		if err != nil {
			return tag, err
		}
		tag.Payload = v
	case TypeList:
		v, err := readList(r)
		if err != nil {
			return tag, err
		}
		tag.Payload = v
	case TypeCompound:
		v, err := readCompound(r)
		if err != nil {
			return tag, err
		}
		tag.Payload = v
	case TypeIntArray:
		v, err := readIntArray(r)
		if err != nil {
			return tag, err
		}
		tag.Payload = v
	case TypeLongArray:
		v, err := readLongArray(r)
		if err != nil {
			return tag, err
		}
		tag.Payload = v
	default:
		return tag, ErrUnknownTag
	}

	return tag, nil
}

func readString(r io.Reader) (String, error) {
	var n Short
	if err := read(r, &n); err != nil {
		return "", err
	}

	if n < 0 {
		return "", ErrNegativeLength
	}

	b := make([]byte, n)
	if err := read(r, b); err != nil {
		return "", err
	}

	return String(b), nil
}

func readByteArray(r io.Reader) (ByteArray, error) {
	var n Int
	if err := read(r, &n); err != nil {
		return nil, err
	}

	if n < 0 {
		return nil, ErrNegativeLength
	}

	a := make(ByteArray, n)
	if err := read(r, a); err != nil {
		return nil, err
	}

	return a, nil
}

func readIntArray(r io.Reader) (IntArray, error) {
	var n Int
	if err := read(r, &n); err != nil {
		return nil, err
	}

	if n < 0 {
		return nil, ErrNegativeLength
	}

	a := make(IntArray, n)
	if err := read(r, a); err != nil {
		return nil, err
	}

	return a, nil
}

func readLongArray(r io.Reader) (LongArray, error) {
	var n Int
	if err := read(r, &n); err != nil {
		return nil, err
	}

	if n < 0 {
		return nil, ErrNegativeLength
	}

	a := make(LongArray, n)
	if err := read(r, a); err != nil {
		return nil, err
	}

	return a, nil
}

func readList(r io.Reader) (List, error) {
	var list List

	if err := read(r, &list.ElementType); err != nil {
		return list, err
	}

	var n Int
	if err := read(r, &n); err != nil {
		return list, err
	}

	if n < 0 {
		return list, ErrNegativeLength
	}

	switch list.ElementType {
	case TypeEnd:
		if n > 0 {
			return list, errors.New("list of end tags has nonzero length")
		}
	case TypeByte:
		a := make([]Byte, n)
		if err := read(r, a); err != nil {
			return list, err
		}
		list.Payload = a
	case TypeShort:
		a := make([]Short, n)
		if err := read(r, a); err != nil {
			return list, err
		}
		list.Payload = a
	case TypeInt:
		a := make([]Int, n)
		if err := read(r, a); err != nil {
			return list, err
		}
		list.Payload = a
	case TypeLong:
		a := make([]Long, n)
		if err := read(r, a); err != nil {
			return list, err
		}
		list.Payload = a
	case TypeFloat:
		a := make([]Float, n)
		if err := read(r, a); err != nil {
			return list, err
		}
		list.Payload = a
	case TypeDouble:
		a := make([]Double, n)
		if err := read(r, a); err != nil {
			return list, err
		}
		list.Payload = a
	case TypeByteArray:
		a := make([]ByteArray, n)
		for i := range a {
			v, err := readByteArray(r)
			if err != nil {
				return list, err
			}
			a[i] = v
		}
		list.Payload = a
	case TypeString:
		a := make([]String, n)
		for i := range a {
			v, err := readString(r)
			if err != nil {
				return list, err
			}
			a[i] = v
		}
		list.Payload = a
	case TypeList:
		a := make([]List, n)
		for i := range a {
			v, err := readList(r)
			if err != nil {
				return list, err
			}
			a[i] = v
		}
		list.Payload = a
	case TypeCompound:
		a := make([]Compound, n)
		for i := range a {
			v, err := readCompound(r)
			if err != nil {
				return list, err
			}
			a[i] = v
		}
		list.Payload = a
	case TypeIntArray:
		a := make([]IntArray, n)
		for i := range a {
			v, err := readIntArray(r)
			if err != nil {
				return list, err
			}
			a[i] = v
		}
		list.Payload = a
	case TypeLongArray:
		a := make([]LongArray, n)
		for i := range a {
			v, err := readLongArray(r)
			if err != nil {
				return list, err
			}
			a[i] = v
		}
		list.Payload = a
	default:
		return list, ErrUnknownTag
	}

	return list, nil
}

func readCompound(r io.Reader) (Compound, error) {
	m := make(Compound)

	for {
		tag, err := readNamedTag(r)
		if err != nil {
			return nil, err
		}

		if tag.Type == TypeEnd {
			break
		}

		if _, exists := m[tag.Name]; exists {
			return nil, ErrDuplicateName
		}
		m[tag.Name] = tag
	}

	return m, nil
}
