package nbt

import (
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
)

type Decoder struct {
	r *offsetReader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: &offsetReader{r: r}}
}

type offsetReader struct {
	r      io.Reader
	offset int64
}

func (r *offsetReader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	r.offset += int64(n)
	return n, err
}

func (dec *Decoder) Decode() (*NamedTag, error) {
	tag, err := dec.readNamedTag()
	if err != nil {
		return nil, errors.WithMessagef(err, "offset %d:", dec.r.offset)
	}
	return tag, nil
}

func readBE(r io.Reader, v interface{}) error {
	return binary.Read(r, binary.BigEndian, v)
}

func (dec *Decoder) readNamedTag() (*NamedTag, error) {
	var typ Type
	if err := readBE(dec.r, &typ); err != nil {
		return nil, err
	}

	if typ == TypeEnd {
		return &NamedTag{}, nil
	}

	name, err := dec.readString()
	if err != nil {
		return nil, err
	}

	var payload interface{}
	switch typ {
	case TypeByte:
		var n int8
		err = readBE(dec.r, &n)
		payload = n
	case TypeShort:
		var n int16
		err = readBE(dec.r, &n)
		payload = n
	case TypeInt:
		var n int32
		err = readBE(dec.r, &n)
		payload = n
	case TypeLong:
		var n int64
		err = readBE(dec.r, &n)
		payload = n
	case TypeFloat:
		var x float32
		err = readBE(dec.r, &x)
		payload = x
	case TypeDouble:
		var x float64
		err = readBE(dec.r, &x)
		payload = x
	case TypeByteArray:
		payload, err = dec.readByteArray()
	case TypeString:
		payload, err = dec.readString()
	case TypeList:
		payload, err = dec.readList()
	case TypeCompound:
		payload, err = dec.readCompound()
	case TypeIntArray:
		payload, err = dec.readIntArray()
	case TypeLongArray:
		payload, err = dec.readLongArray()
	default:
		return nil, errors.Errorf("unknown type %#x", typ)
	}

	if err != nil {
		return nil, err
	}

	return &NamedTag{typ, name, payload}, nil
}

func (dec *Decoder) readByteArray() ([]byte, error) {
	var length int32
	if err := readBE(dec.r, &length); err != nil {
		return nil, err
	}

	if length < 0 {
		return nil, errors.Errorf("negative length (%d)", length)
	}

	b := make([]byte, length)
	if err := readBE(dec.r, b); err != nil {
		return nil, err
	}

	return b, nil
}

func (dec *Decoder) readString() (string, error) {
	var length int16
	if err := readBE(dec.r, &length); err != nil {
		return "", err
	}

	if length < 0 {
		return "", errors.Errorf("negative length (%d)", length)
	}

	b := make([]byte, length)
	if err := readBE(dec.r, b); err != nil {
		return "", err
	}

	return string(b), nil
}

func (dec *Decoder) readList() (*List, error) {
	var typ Type
	err := readBE(dec.r, &typ)
	if err != nil {
		return nil, err
	}

	var length int32
	if err = readBE(dec.r, &length); err != nil {
		return nil, err
	}

	if length < 0 {
		return nil, errors.Errorf("negative length (%d)", length)
	}

	var array interface{}
	switch typ {
	case TypeEnd:
	case TypeByte:
		a := make([]int8, length)
		if err = readBE(dec.r, a); err != nil {
			return nil, err
		}
		array = a
	case TypeShort:
		a := make([]int16, length)
		if err = readBE(dec.r, a); err != nil {
			return nil, err
		}
		array = a
	case TypeInt:
		a := make([]int32, length)
		if err = readBE(dec.r, a); err != nil {
			return nil, err
		}
		array = a
	case TypeLong:
		a := make([]int64, length)
		if err = readBE(dec.r, a); err != nil {
			return nil, err
		}
		array = a
	case TypeFloat:
		a := make([]float32, length)
		if err = readBE(dec.r, a); err != nil {
			return nil, err
		}
		array = a
	case TypeDouble:
		a := make([]float64, length)
		if err = readBE(dec.r, a); err != nil {
			return nil, err
		}
		array = a
	case TypeByteArray:
		a := make([][]byte, length)
		for i := range a {
			if a[i], err = dec.readByteArray(); err != nil {
				return nil, err
			}
		}
		array = a
	case TypeString:
		a := make([]string, length)
		for i := range a {
			if a[i], err = dec.readString(); err != nil {
				return nil, err
			}
		}
		array = a
	case TypeList:
		a := make([]*List, length)
		for i := range a {
			if a[i], err = dec.readList(); err != nil {
				return nil, err
			}
		}
		array = a
	case TypeCompound:
		a := make([]Compound, length)
		for i := range a {
			if a[i], err = dec.readCompound(); err != nil {
				return nil, err
			}
		}
		array = a
	case TypeIntArray:
		a := make([][]int32, length)
		for i := range a {
			if a[i], err = dec.readIntArray(); err != nil {
				return nil, err
			}
		}
		array = a
	case TypeLongArray:
		a := make([][]int64, length)
		for i := range a {
			if a[i], err = dec.readLongArray(); err != nil {
				return nil, err
			}
		}
		array = a
	default:
		return nil, errors.Errorf("unknown type %#x", typ)
	}

	return &List{typ, array}, nil
}

func (dec *Decoder) readCompound() (Compound, error) {
	m := make(Compound)
	for {
		tag, err := dec.readNamedTag()
		if err != nil {
			return nil, err
		}

		if tag.Type == TypeEnd {
			return m, nil
		}

		if _, exists := m[tag.Name]; exists {
			return nil, errors.Errorf("duplicate name %q", tag.Name)
		}
		m[tag.Name] = tag
	}
}

func (dec *Decoder) readIntArray() ([]int32, error) {
	var length int32
	if err := readBE(dec.r, &length); err != nil {
		return nil, err
	}

	if length < 0 {
		return nil, errors.Errorf("negative length (%d)", length)
	}

	a := make([]int32, length)
	if err := readBE(dec.r, a); err != nil {
		return nil, err
	}

	return a, nil
}

func (dec *Decoder) readLongArray() ([]int64, error) {
	var length int32
	if err := readBE(dec.r, &length); err != nil {
		return nil, err
	}

	if length < 0 {
		return nil, errors.Errorf("negative length (%d)", length)
	}

	a := make([]int64, length)
	if err := readBE(dec.r, a); err != nil {
		return nil, err
	}

	return a, nil
}