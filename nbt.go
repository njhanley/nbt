package nbt

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"runtime"
	"strconv"
)

func read(r io.Reader, v interface{}) error {
	return binary.Read(r, binary.BigEndian, v)
}

func write(w io.Writer, v interface{}) error {
	return binary.Write(w, binary.BigEndian, v)
}

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

var (
	typeNames = []string{
		TypeEnd:       "end",
		TypeByte:      "byte",
		TypeShort:     "short",
		TypeInt:       "int",
		TypeLong:      "long",
		TypeFloat:     "float",
		TypeDouble:    "double",
		TypeByteArray: "byteArray",
		TypeString:    "string",
		TypeList:      "list",
		TypeCompound:  "compound",
		TypeIntArray:  "intArray",
		TypeLongArray: "longArray",
	}
	typeIDs map[string]Type
)

func init() {
	// create map from type names to IDs
	typeIDs = make(map[string]Type, len(typeNames))
	for id, name := range typeNames {
		typeIDs[name] = Type(id)
	}
}

func (typ Type) String() string {
	return typeNames[typ]
}

func (typ Type) MarshalJSON() ([]byte, error) {
	return json.Marshal(typ.String())
}

func (typ *Type) UnmarshalJSON(data []byte) error {
	s, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	id, ok := typeIDs[s]
	if !ok {
		return fmt.Errorf("unknown tag type %q", s)
	}
	*typ = id
	return nil
}

func (typ Type) Decode(r io.Reader) (interface{}, error) {
	switch typ {
	case TypeByte:
		return readByte(r)
	case TypeShort:
		return readShort(r)
	case TypeInt:
		return readInt(r)
	case TypeLong:
		return readLong(r)
	case TypeFloat:
		return readFloat(r)
	case TypeDouble:
		return readDouble(r)
	case TypeByteArray:
		return readByteArray(r)
	case TypeString:
		return readString(r)
	case TypeList:
		return readList(r)
	case TypeCompound:
		return readCompound(r)
	case TypeIntArray:
		return readIntArray(r)
	case TypeLongArray:
		return readLongArray(r)
	default:
		return nil, errors.New("no decoder for type")
	}
}

func (typ Type) Encode(w io.Writer, v interface{}) (err error) {
	// rather than explicitly checking each type assertion,
	// recover and return the TypeAssertionError if one is raised
	defer func() {
		if x := recover(); x != nil {
			if _err, ok := x.(*runtime.TypeAssertionError); ok {
				err = _err
				return
			}
			panic(x)
		}
	}()

	switch typ {
	case TypeByte:
		return writeByte(w, v.(Byte))
	case TypeShort:
		return writeShort(w, v.(Short))
	case TypeInt:
		return writeInt(w, v.(Int))
	case TypeLong:
		return writeLong(w, v.(Long))
	case TypeFloat:
		return writeFloat(w, v.(Float))
	case TypeDouble:
		return writeDouble(w, v.(Double))
	case TypeByteArray:
		return writeByteArray(w, v.(ByteArray))
	case TypeString:
		return writeString(w, v.(String))
	case TypeList:
		return writeList(w, v.(List))
	case TypeCompound:
		return writeCompound(w, v.(Compound))
	case TypeIntArray:
		return writeIntArray(w, v.(IntArray))
	case TypeLongArray:
		return writeLongArray(w, v.(LongArray))
	default:
		return errors.New("no encoder for type")
	}
}

func (typ Type) UnmarshalJSONPayload(data []byte) (interface{}, error) {
	switch typ {
	case TypeByte:
		x := new(Byte)
		err := json.Unmarshal(data, x)
		return *x, err
	case TypeShort:
		x := new(Short)
		err := json.Unmarshal(data, x)
		return *x, err
	case TypeInt:
		x := new(Int)
		err := json.Unmarshal(data, x)
		return *x, err
	case TypeLong:
		x := new(Long)
		err := json.Unmarshal(data, x)
		return *x, err
	case TypeFloat:
		x := new(Float)
		err := json.Unmarshal(data, x)
		return *x, err
	case TypeDouble:
		x := new(Double)
		err := json.Unmarshal(data, x)
		return *x, err
	case TypeByteArray:
		x := new(ByteArray)
		err := json.Unmarshal(data, x)
		return *x, err
	case TypeString:
		x := new(String)
		err := json.Unmarshal(data, x)
		return *x, err
	case TypeList:
		x := new(List)
		err := json.Unmarshal(data, x)
		return *x, err
	case TypeCompound:
		x := new(Compound)
		err := json.Unmarshal(data, x)
		return *x, err
	case TypeIntArray:
		x := new(IntArray)
		err := json.Unmarshal(data, x)
		return *x, err
	case TypeLongArray:
		x := new(LongArray)
		err := json.Unmarshal(data, x)
		return *x, err
	default:
		return nil, errors.New("invalid type")
	}
}

func readType(r io.Reader) (Type, error) {
	var typ Type
	err := read(r, &typ)
	return typ, err
}

func writeType(w io.Writer, typ Type) error {
	return write(w, typ)
}

type Tag struct {
	Type    Type
	Payload interface{}
}

func (tag *Tag) UnmarshalJSON(data []byte) error {
	var x struct {
		Type    Type
		Payload json.RawMessage
	}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	tag.Type = x.Type

	v, err := tag.Type.UnmarshalJSONPayload(x.Payload)
	if err != nil {
		return err
	}
	tag.Payload = v

	return nil
}

type NamedTag struct {
	Tag
	Name String
}

func (tag *NamedTag) UnmarshalJSON(data []byte) error {
	var x struct {
		Type    Type
		Name    String
		Payload json.RawMessage
	}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	tag.Type = x.Type
	tag.Name = x.Name

	v, err := tag.Type.UnmarshalJSONPayload(x.Payload)
	if err != nil {
		return err
	}
	tag.Payload = v

	return nil
}

func readNamedTag(r io.Reader) (NamedTag, error) {
	var tag NamedTag

	typ, err := readType(r)
	if err != nil {
		return tag, err
	}
	tag.Type = typ

	if tag.Type == TypeEnd {
		return tag, nil
	}

	tag.Name, err = readString(r)
	if err != nil {
		return tag, err
	}

	v, err := tag.Type.Decode(r)
	if err != nil {
		return tag, err
	}
	tag.Payload = v

	return tag, nil
}

func writeNamedTag(w io.Writer, tag NamedTag) error {
	if err := writeType(w, tag.Type); err != nil {
		return err
	}

	if tag.Type == TypeEnd {
		return nil
	}

	if err := writeString(w, tag.Name); err != nil {
		return err
	}

	if err := tag.Type.Encode(w, tag.Payload); err != nil {
		return err
	}

	return nil
}

func Decode(r io.Reader) (NamedTag, error) {
	return readNamedTag(r)
}

func Encode(w io.Writer, tag NamedTag) error {
	return writeNamedTag(w, tag)
}

const (
	MaxByte  = 1<<7 - 1
	MaxShort = 1<<15 - 1
	MaxInt   = 1<<31 - 1
	MaxLong  = 1<<63 - 1
)

type Byte int8

func (n Byte) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(strconv.FormatInt(int64(n), 10))), nil
}

func (n *Byte) UnmarshalJSON(data []byte) error {
	s, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	_n, err := strconv.ParseInt(s, 10, 8)
	*n = Byte(_n)
	return err
}

func readByte(r io.Reader) (Byte, error) {
	var n Byte
	err := read(r, &n)
	return n, err
}

func writeByte(w io.Writer, n Byte) error {
	return write(w, n)
}

type Short int16

func (n Short) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(strconv.FormatInt(int64(n), 10))), nil
}

func (n *Short) UnmarshalJSON(data []byte) error {
	s, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	_n, err := strconv.ParseInt(s, 10, 16)
	*n = Short(_n)
	return err
}

func readShort(r io.Reader) (Short, error) {
	var n Short
	err := read(r, &n)
	return n, err
}

func writeShort(w io.Writer, n Short) error {
	return write(w, n)
}

type Int int32

func (n Int) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(strconv.FormatInt(int64(n), 10))), nil
}

func (n *Int) UnmarshalJSON(data []byte) error {
	s, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	_n, err := strconv.ParseInt(s, 10, 32)
	*n = Int(_n)
	return err
}

func readInt(r io.Reader) (Int, error) {
	var n Int
	err := read(r, &n)
	return n, err
}

func writeInt(w io.Writer, n Int) error {
	return write(w, n)
}

type Long int64

func (n Long) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(strconv.FormatInt(int64(n), 10))), nil
}

func (n *Long) UnmarshalJSON(data []byte) error {
	s, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	_n, err := strconv.ParseInt(s, 10, 64)
	*n = Long(_n)
	return err
}

func readLong(r io.Reader) (Long, error) {
	var n Long
	err := read(r, &n)
	return n, err
}

func writeLong(w io.Writer, n Long) error {
	return write(w, n)
}

type Float float32

func (f Float) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(strconv.FormatFloat(float64(f), 'g', -1, 32))), nil
}

func (f *Float) UnmarshalJSON(data []byte) error {
	s, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	_f, err := strconv.ParseFloat(s, 32)
	*f = Float(_f)
	return err
}

func readFloat(r io.Reader) (Float, error) {
	var f Float
	err := read(r, &f)
	return f, err
}

func writeFloat(w io.Writer, f Float) error {
	return write(w, f)
}

type Double float64

func (f Double) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(strconv.FormatFloat(float64(f), 'g', -1, 64))), nil
}

func (f *Double) UnmarshalJSON(data []byte) error {
	s, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	_f, err := strconv.ParseFloat(s, 64)
	*f = Double(_f)
	return err
}

func readDouble(r io.Reader) (Double, error) {
	var f Double
	err := read(r, &f)
	return f, err
}

func writeDouble(w io.Writer, f Double) error {
	return write(w, f)
}

type ByteArray []Byte

func readByteArray(r io.Reader) (ByteArray, error) {
	n, err := readInt(r)
	if err != nil {
		return nil, err
	}

	if n < 0 {
		return nil, errors.New("negative length")
	}

	a := make(ByteArray, n)
	if err := read(r, a); err != nil {
		return nil, err
	}

	return a, nil
}

func writeByteArray(w io.Writer, a ByteArray) error {
	n := len(a)
	if n > MaxInt {
		return fmt.Errorf("byteArray length (%d) > maxInt (%d)", n, MaxInt)
	}

	if err := writeInt(w, Int(n)); err != nil {
		return err
	}

	return write(w, a)
}

type String string

func readString(r io.Reader) (String, error) {
	n, err := readShort(r)
	if err != nil {
		return "", err
	}

	if n < 0 {
		return "", errors.New("negative length")
	}

	b := make([]byte, n)
	if err := read(r, b); err != nil {
		return "", err
	}

	return String(b), nil
}

func writeString(w io.Writer, s String) error {
	n := len(s)
	if n > MaxShort {
		return fmt.Errorf("string length (%d) > maxShort (%d)", n, MaxShort)
	}

	if err := writeShort(w, Short(n)); err != nil {
		return err
	}

	return write(w, []byte(s))
}

type List struct {
	Type     Type
	Elements []interface{}
}

func (l *List) UnmarshalJSON(data []byte) error {
	var x struct {
		Type     Type
		Elements []json.RawMessage
	}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	l.Type = x.Type

	l.Elements = make([]interface{}, len(x.Elements))
	for i, m := range x.Elements {
		v, err := l.Type.UnmarshalJSONPayload(m)
		if err != nil {
			return err
		}
		l.Elements[i] = v
	}

	return nil
}

func readList(r io.Reader) (List, error) {
	var l List

	typ, err := readType(r)
	if err != nil {
		return l, err
	}
	if typ == TypeEnd {
		return l, errors.New("invalid list type")
	}
	l.Type = typ

	n, err := readInt(r)
	if err != nil {
		return l, err
	}
	if n < 0 {
		return l, errors.New("negative length")
	}

	a := make([]interface{}, int(n))
	for i := range a {
		v, err := l.Type.Decode(r)
		if err != nil {
			return l, err
		}
		a[i] = v
	}
	l.Elements = a

	return l, nil
}

func writeList(w io.Writer, l List) error {
	if l.Type == TypeEnd {
		return errors.New("invalid list type")
	}

	if err := writeType(w, l.Type); err != nil {
		return err
	}

	n := len(l.Elements)
	if n > MaxInt {
		return errors.New("list is too long")
	}
	if err := writeInt(w, Int(n)); err != nil {
		return err
	}

	for _, v := range l.Elements {
		if err := l.Type.Encode(w, v); err != nil {
			return err
		}
	}

	return nil
}

type Compound map[String]Tag

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
			return nil, errors.New("duplicate name")
		}
		m[tag.Name] = tag.Tag
	}
	return m, nil
}

func writeCompound(w io.Writer, m Compound) error {
	for name, tag := range m {
		if err := writeNamedTag(w, NamedTag{Tag: tag, Name: name}); err != nil {
			return err
		}
	}
	return writeNamedTag(w, NamedTag{Tag: Tag{Type: TypeEnd}})
}

type IntArray []Int

func readIntArray(r io.Reader) (IntArray, error) {
	n, err := readInt(r)
	if err != nil {
		return nil, err
	}

	if n < 0 {
		return nil, errors.New("negative length")
	}

	a := make(IntArray, n)
	for i := range a {
		a[i], err = readInt(r)
		if err != nil {
			return nil, err
		}
	}

	return a, nil
}

func writeIntArray(w io.Writer, a IntArray) error {
	n := len(a)
	if n > MaxInt {
		return fmt.Errorf("intArray length (%d) > maxInt (%d)", n, MaxInt)
	}

	if err := writeInt(w, Int(n)); err != nil {
		return err
	}

	for _, v := range a {
		if err := writeInt(w, v); err != nil {
			return err
		}
	}

	return nil
}

type LongArray []Long

func readLongArray(r io.Reader) (LongArray, error) {
	n, err := readInt(r)
	if err != nil {
		return nil, err
	}

	if n < 0 {
		return nil, errors.New("negative length")
	}

	a := make(LongArray, n)
	for i := range a {
		a[i], err = readLong(r)
		if err != nil {
			return nil, err
		}
	}

	return a, nil
}

func writeLongArray(w io.Writer, a LongArray) error {
	n := len(a)
	if n > MaxInt {
		return fmt.Errorf("longArray length (%d) > maxInt (%d)", n, MaxInt)
	}

	if err := writeInt(w, Int(n)); err != nil {
		return err
	}

	for _, v := range a {
		if err := writeLong(w, v); err != nil {
			return err
		}
	}

	return nil
}
