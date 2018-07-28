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

type NBT struct {
	Name string
	Tags map[string]interface{}
}

func Decode(r io.Reader) (*NBT, error) {
	var id byte
	if err := read(r, &id); err != nil {
		return nil, err
	}

	if id != idCompound {
		return nil, ErrInvalidRoot
	}

	name, err := readString(r)
	if err != nil {
		return nil, err
	}

	tags, err := readCompound(r)
	if err != nil {
		return nil, err
	}

	return &NBT{name, tags}, nil
}

func read(r io.Reader, v interface{}) error {
	return binary.Read(r, binary.BigEndian, v)
}

func readString(r io.Reader) (string, error) {
	var n int16
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

	return string(b), nil
}

func readByteArray(r io.Reader) ([]int8, error) {
	var n int32
	if err := read(r, &n); err != nil {
		return nil, err
	}

	if n < 0 {
		return nil, ErrNegativeLength
	}

	a := make([]int8, n)
	if err := read(r, a); err != nil {
		return nil, err
	}

	return a, nil
}

func readIntArray(r io.Reader) ([]int32, error) {
	var n int32
	if err := read(r, &n); err != nil {
		return nil, err
	}

	if n < 0 {
		return nil, ErrNegativeLength
	}

	a := make([]int32, n)
	if err := read(r, a); err != nil {
		return nil, err
	}

	return a, nil
}

func readLongArray(r io.Reader) ([]int64, error) {
	var n int32
	if err := read(r, &n); err != nil {
		return nil, err
	}

	if n < 0 {
		return nil, ErrNegativeLength
	}

	a := make([]int64, n)
	if err := read(r, a); err != nil {
		return nil, err
	}

	return a, nil
}

const (
	idEnd byte = iota
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

func readList(r io.Reader) (interface{}, error) {
	var id byte
	if err := read(r, &id); err != nil {
		return nil, err
	}

	var n int32
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
		return []struct{}{}, nil
	case idByte:
		a := make([]int8, n)
		if err := read(r, a); err != nil {
			return nil, err
		}
		return a, nil
	case idShort:
		a := make([]int16, n)
		if err := read(r, a); err != nil {
			return nil, err
		}
		return a, nil
	case idInt:
		a := make([]int32, n)
		if err := read(r, a); err != nil {
			return nil, err
		}
		return a, nil
	case idLong:
		a := make([]int64, n)
		if err := read(r, a); err != nil {
			return nil, err
		}
		return a, nil
	case idFloat:
		a := make([]float32, n)
		if err := read(r, a); err != nil {
			return nil, err
		}
		return a, nil
	case idDouble:
		a := make([]float64, n)
		if err := read(r, a); err != nil {
			return nil, err
		}
		return a, nil
	case idByteArray:
		a := make([][]int8, n)
		for i := range a {
			v, err := readByteArray(r)
			if err != nil {
				return nil, err
			}
			a[i] = v
		}
		return a, nil
	case idString:
		a := make([]string, n)
		for i := range a {
			v, err := readString(r)
			if err != nil {
				return nil, err
			}
			a[i] = v
		}
		return a, nil
	case idList:
		a := make([]interface{}, n)
		for i := range a {
			v, err := readList(r)
			if err != nil {
				return nil, err
			}
			a[i] = v
		}
		return a, nil
	case idCompound:
		a := make([]map[string]interface{}, n)
		for i := range a {
			v, err := readCompound(r)
			if err != nil {
				return nil, err
			}
			a[i] = v
		}
		return a, nil
	case idIntArray:
		a := make([][]int32, n)
		for i := range a {
			v, err := readIntArray(r)
			if err != nil {
				return nil, err
			}
			a[i] = v
		}
		return a, nil
	case idLongArray:
		a := make([][]int64, n)
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

func readCompound(r io.Reader) (map[string]interface{}, error) {
	m := make(map[string]interface{})

	for {
		var id byte
		if err := read(r, &id); err != nil {
			return nil, err
		}

		if id == idEnd {
			break
		}

		name, err := readString(r)
		if err != nil {
			return nil, err
		}

		if _, exists := m[name]; exists {
			return nil, ErrDuplicateName
		}

		switch id {
		case idByte:
			var v int8
			if err := read(r, &v); err != nil {
				return nil, err
			}
			m[name] = v
		case idShort:
			var v int16
			if err := read(r, &v); err != nil {
				return nil, err
			}
			m[name] = v
		case idInt:
			var v int32
			if err := read(r, &v); err != nil {
				return nil, err
			}
			m[name] = v
		case idLong:
			var v int64
			if err := read(r, &v); err != nil {
				return nil, err
			}
			m[name] = v
		case idFloat:
			var v float32
			if err := read(r, &v); err != nil {
				return nil, err
			}
			m[name] = v
		case idDouble:
			var v float64
			if err := read(r, &v); err != nil {
				return nil, err
			}
			m[name] = v
		case idByteArray:
			v, err := readByteArray(r)
			if err != nil {
				return nil, err
			}
			m[name] = v
		case idString:
			v, err := readString(r)
			if err != nil {
				return nil, err
			}
			m[name] = v
		case idList:
			v, err := readList(r)
			if err != nil {
				return nil, err
			}
			m[name] = v
		case idCompound:
			v, err := readCompound(r)
			if err != nil {
				return nil, err
			}
			m[name] = v
		case idIntArray:
			v, err := readIntArray(r)
			if err != nil {
				return nil, err
			}
			m[name] = v
		case idLongArray:
			v, err := readLongArray(r)
			if err != nil {
				return nil, err
			}
			m[name] = v
		default:
			return nil, ErrUnknownTag
		}
	}

	return m, nil
}
