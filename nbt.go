package nbt

import (
	"encoding/binary"
	"errors"
	"io"
)

var (
	ErrNegativeLength = errors.New("length is negative")
	ErrUnknownTag     = errors.New("unknown tag ID")
	ErrInvalidRoot    = errors.New("root tag is not compound")
	ErrDuplicateName  = errors.New("duplicate name in compound")
)

const (
	idEnd Byte = iota
	idByte
	idShort
	idInt
	idLong
	idFloat
	idDouble
	idByteArray
	idString
	idList
	idCompound
	idIntArray
	idLongArray
)

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
	List      interface{}
	Compound  map[String]interface{}
	IntArray  []Int
	LongArray []Long
)

func read(r io.Reader, v interface{}) error {
	return binary.Read(r, binary.BigEndian, v)
}

func Decode(r io.Reader) (NamedTag, error) {
	return readNamedTag(r)
}

type NamedTag struct {
	ID      Byte
	Name    String
	Payload interface{}
}

func readNamedTag(r io.Reader) (NamedTag, error) {
	var tag NamedTag

	err := read(r, &tag.ID)
	if err != nil {
		return tag, err
	}

	if tag.ID == idEnd {
		return tag, nil
	}

	tag.Name, err = readString(r)
	if err != nil {
		return tag, err
	}

	switch tag.ID {
	case idByte:
		var v Byte
		if err := read(r, &v); err != nil {
			return tag, err
		}
		tag.Payload = v
	case idShort:
		var v Short
		if err := read(r, &v); err != nil {
			return tag, err
		}
		tag.Payload = v
	case idInt:
		var v Int
		if err := read(r, &v); err != nil {
			return tag, err
		}
		tag.Payload = v
	case idLong:
		var v Long
		if err := read(r, &v); err != nil {
			return tag, err
		}
		tag.Payload = v
	case idFloat:
		var v Float
		if err := read(r, &v); err != nil {
			return tag, err
		}
		tag.Payload = v
	case idDouble:
		var v Double
		if err := read(r, &v); err != nil {
			return tag, err
		}
		tag.Payload = v
	case idByteArray:
		v, err := readByteArray(r)
		if err != nil {
			return tag, err
		}
		tag.Payload = v
	case idString:
		v, err := readString(r)
		if err != nil {
			return tag, err
		}
		tag.Payload = v
	case idList:
		v, err := readList(r)
		if err != nil {
			return tag, err
		}
		tag.Payload = v
	case idCompound:
		v, err := readCompound(r)
		if err != nil {
			return tag, err
		}
		tag.Payload = v
	case idIntArray:
		v, err := readIntArray(r)
		if err != nil {
			return tag, err
		}
		tag.Payload = v
	case idLongArray:
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

func readList(r io.Reader) (interface{}, error) {
	var id Byte
	if err := read(r, &id); err != nil {
		return nil, err
	}

	var n Int
	if err := read(r, &n); err != nil {
		return nil, err
	}

	if n < 0 {
		return nil, ErrNegativeLength
	}

	switch id {
	case idEnd:
		if n > 0 {
			return nil, errors.New("list of end tags has nonzero length")
		}
		return []End{}, nil
	case idByte:
		a := make([]Byte, n)
		if err := read(r, a); err != nil {
			return nil, err
		}
		return a, nil
	case idShort:
		a := make([]Short, n)
		if err := read(r, a); err != nil {
			return nil, err
		}
		return a, nil
	case idInt:
		a := make([]Int, n)
		if err := read(r, a); err != nil {
			return nil, err
		}
		return a, nil
	case idLong:
		a := make([]Long, n)
		if err := read(r, a); err != nil {
			return nil, err
		}
		return a, nil
	case idFloat:
		a := make([]Float, n)
		if err := read(r, a); err != nil {
			return nil, err
		}
		return a, nil
	case idDouble:
		a := make([]Double, n)
		if err := read(r, a); err != nil {
			return nil, err
		}
		return a, nil
	case idByteArray:
		a := make([]ByteArray, n)
		for i := range a {
			v, err := readByteArray(r)
			if err != nil {
				return nil, err
			}
			a[i] = v
		}
		return a, nil
	case idString:
		a := make([]String, n)
		for i := range a {
			v, err := readString(r)
			if err != nil {
				return nil, err
			}
			a[i] = v
		}
		return a, nil
	case idList:
		a := make([]List, n)
		for i := range a {
			v, err := readList(r)
			if err != nil {
				return nil, err
			}
			a[i] = v
		}
		return a, nil
	case idCompound:
		a := make([]Compound, n)
		for i := range a {
			v, err := readCompound(r)
			if err != nil {
				return nil, err
			}
			a[i] = v
		}
		return a, nil
	case idIntArray:
		a := make([]IntArray, n)
		for i := range a {
			v, err := readIntArray(r)
			if err != nil {
				return nil, err
			}
			a[i] = v
		}
		return a, nil
	case idLongArray:
		a := make([]LongArray, n)
		for i := range a {
			v, err := readLongArray(r)
			if err != nil {
				return nil, err
			}
			a[i] = v
		}
		return a, nil
	default:
		return nil, ErrUnknownTag
	}
}

func readCompound(r io.Reader) (Compound, error) {
	m := make(Compound)

	for {
		tag, err := readNamedTag(r)
		if err != nil {
			return nil, err
		}

		if tag.ID == idEnd {
			break
		}

		if _, exists := m[tag.Name]; exists {
			return nil, ErrDuplicateName
		}
		m[tag.Name] = tag.Payload
	}

	return m, nil
}
