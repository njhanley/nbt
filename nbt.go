package nbt

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
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

const (
	maxByte  = 1<<7 - 1
	maxShort = 1<<15 - 1
	maxInt   = 1<<31 - 1
	maxLong  = 1<<63 - 1
)

func (m Compound) MarshalJSON() ([]byte, error) {
	a := make([]NamedTag, 0, len(m))
	for _, tag := range m {
		a = append(a, tag)
	}
	sort.Slice(a, func(i, j int) bool { return a[i].Name < a[j].Name })
	return json.Marshal(a)
}

func read(r io.Reader, v interface{}) error {
	return binary.Read(r, binary.BigEndian, v)
}

func write(w io.Writer, v interface{}) error {
	return binary.Write(w, binary.BigEndian, v)
}

func Decode(r io.Reader) (NamedTag, error) {
	return readNamedTag(r)
}

func Encode(w io.Writer, tag NamedTag) error {
	return writeNamedTag(w, tag)
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

func writeNamedTag(w io.Writer, tag NamedTag) error {
	if err := write(w, tag.Type); err != nil {
		return err
	}

	if tag.Type == TypeEnd {
		return nil
	}

	if err := writeString(w, tag.Name); err != nil {
		return err
	}

	switch tag.Type {
	case TypeByte, TypeShort, TypeInt, TypeLong, TypeFloat, TypeDouble:
		if err := write(w, tag.Payload); err != nil {
			return err
		}
	case TypeByteArray:
		if err := writeByteArray(w, tag.Payload.(ByteArray)); err != nil {
			return err
		}
	case TypeString:
		if err := writeString(w, tag.Payload.(String)); err != nil {
			return err
		}
	case TypeList:
		if err := writeList(w, tag.Payload.(List)); err != nil {
			return err
		}
	case TypeCompound:
		if err := writeCompound(w, tag.Payload.(Compound)); err != nil {
			return err
		}
	case TypeIntArray:
		if err := writeIntArray(w, tag.Payload.(IntArray)); err != nil {
			return err
		}
	case TypeLongArray:
		if err := writeLongArray(w, tag.Payload.(LongArray)); err != nil {
			return err
		}
	default:
		return ErrUnknownTag
	}

	return nil
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

func writeString(w io.Writer, s String) error {
	n := len(s)
	if n > maxShort {
		return fmt.Errorf("string length (%d) > maxShort (%d)", n, maxShort)
	}

	if err := write(w, Short(n)); err != nil {
		return err
	}

	return write(w, []byte(s))
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

func writeByteArray(w io.Writer, a ByteArray) error {
	n := len(a)
	if n > maxInt {
		return fmt.Errorf("byteArray length (%d) > maxInt (%d)", n, maxInt)
	}

	if err := write(w, Int(n)); err != nil {
		return err
	}

	return write(w, a)
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

func writeIntArray(w io.Writer, a IntArray) error {
	n := len(a)
	if n > maxInt {
		return fmt.Errorf("intArray length (%d) > maxInt (%d)", n, maxInt)
	}

	if err := write(w, Int(n)); err != nil {
		return err
	}

	return write(w, a)
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

func writeLongArray(w io.Writer, a LongArray) error {
	n := len(a)
	if n > maxInt {
		return fmt.Errorf("longArray length (%d) > maxInt (%d)", n, maxInt)
	}

	if err := write(w, Int(n)); err != nil {
		return err
	}

	return write(w, a)
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

func writeList(w io.Writer, list List) error {
	if err := write(w, list.ElementType); err != nil {
		return err
	}

	switch list.ElementType {
	case TypeEnd:
		if list.Payload != nil {
			return fmt.Errorf("list of end tags has nonempty payload")
		}

		if err := write(w, Int(0)); err != nil {
			return err
		}
	case TypeByte:
		a := list.Payload.([]Byte)

		n := len(a)
		if n > maxInt {
			return fmt.Errorf("list length (%d) > maxInt (%d)", n, maxInt)
		}

		if err := write(w, Int(n)); err != nil {
			return err
		}

		if err := write(w, a); err != nil {
			return err
		}
	case TypeShort:
		a := list.Payload.([]Short)

		n := len(a)
		if n > maxInt {
			return fmt.Errorf("list length (%d) > maxInt (%d)", n, maxInt)
		}

		if err := write(w, Int(n)); err != nil {
			return err
		}

		if err := write(w, a); err != nil {
			return err
		}
	case TypeInt:
		a := list.Payload.([]Int)

		n := len(a)
		if n > maxInt {
			return fmt.Errorf("list length (%d) > maxInt (%d)", n, maxInt)
		}

		if err := write(w, Int(n)); err != nil {
			return err
		}

		if err := write(w, a); err != nil {
			return err
		}
	case TypeLong:
		a := list.Payload.([]Long)

		n := len(a)
		if n > maxInt {
			return fmt.Errorf("list length (%d) > maxInt (%d)", n, maxInt)
		}

		if err := write(w, Int(n)); err != nil {
			return err
		}

		if err := write(w, a); err != nil {
			return err
		}
	case TypeFloat:
		a := list.Payload.([]Float)

		n := len(a)
		if n > maxInt {
			return fmt.Errorf("list length (%d) > maxInt (%d)", n, maxInt)
		}

		if err := write(w, Int(n)); err != nil {
			return err
		}

		if err := write(w, a); err != nil {
			return err
		}
	case TypeDouble:
		a := list.Payload.([]Double)

		n := len(a)
		if n > maxInt {
			return fmt.Errorf("list length (%d) > maxInt (%d)", n, maxInt)
		}

		if err := write(w, Int(n)); err != nil {
			return err
		}

		if err := write(w, a); err != nil {
			return err
		}
	case TypeByteArray:
		a := list.Payload.([]ByteArray)

		n := len(a)
		if n > maxInt {
			return fmt.Errorf("list length (%d) > maxInt (%d)", n, maxInt)
		}

		if err := write(w, Int(n)); err != nil {
			return err
		}

		for i := range a {
			if err := writeByteArray(w, a[i]); err != nil {
				return err
			}
		}
	case TypeString:
		a := list.Payload.([]String)

		n := len(a)
		if n > maxInt {
			return fmt.Errorf("list length (%d) > maxInt (%d)", n, maxInt)
		}

		if err := write(w, Int(n)); err != nil {
			return err
		}

		for i := range a {
			if err := writeString(w, a[i]); err != nil {
				return err
			}
		}
	case TypeList:
		a := list.Payload.([]List)

		n := len(a)
		if n > maxInt {
			return fmt.Errorf("list length (%d) > maxInt (%d)", n, maxInt)
		}

		if err := write(w, Int(n)); err != nil {
			return err
		}

		for i := range a {
			if err := writeList(w, a[i]); err != nil {
				return err
			}
		}
	case TypeCompound:
		a := list.Payload.([]Compound)

		n := len(a)
		if n > maxInt {
			return fmt.Errorf("list length (%d) > maxInt (%d)", n, maxInt)
		}

		if err := write(w, Int(n)); err != nil {
			return err
		}

		for i := range a {
			if err := writeCompound(w, a[i]); err != nil {
				return err
			}
		}
	case TypeIntArray:
		a := list.Payload.([]IntArray)

		n := len(a)
		if n > maxInt {
			return fmt.Errorf("list length (%d) > maxInt (%d)", n, maxInt)
		}

		if err := write(w, Int(n)); err != nil {
			return err
		}

		for i := range a {
			if err := writeIntArray(w, a[i]); err != nil {
				return err
			}
		}
	case TypeLongArray:
		a := list.Payload.([]LongArray)

		n := len(a)
		if n > maxInt {
			return fmt.Errorf("list length (%d) > maxInt (%d)", n, maxInt)
		}

		if err := write(w, Int(n)); err != nil {
			return err
		}

		for i := range a {
			if err := writeLongArray(w, a[i]); err != nil {
				return err
			}
		}
	default:
		return ErrUnknownTag
	}

	return nil
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

func writeCompound(w io.Writer, m Compound) error {
	for _, tag := range m {
		if err := writeNamedTag(w, tag); err != nil {
			return err
		}
	}

	return write(w, TypeEnd)
}
