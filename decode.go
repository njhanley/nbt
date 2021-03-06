package nbt

import (
	"encoding/binary"
	"fmt"
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
	return dec.readNamedTag()
}

func (dec *Decoder) wrap(err error) error {
	if err != nil {
		return &DecodeError{dec.r.offset, errors.WithStack(err)}
	}
	return nil
}

func (dec *Decoder) errorf(format string, a ...interface{}) error {
	return dec.wrap(fmt.Errorf(format, a...))
}

type DecodeError struct {
	Offset int64
	Err    error
}

func (e *DecodeError) Error() string {
	return e.Err.Error()
}

func (e *DecodeError) Format(f fmt.State, c rune) {
	if f.Flag('+') {
		fmt.Fprintf(f, "offset %d: %+v", e.Offset, e.Err)
	} else {
		fmt.Fprint(f, e.Err)
	}
}

func (e *DecodeError) Cause() error {
	return e.Err
}

func readBE(r io.Reader, v interface{}) error {
	return binary.Read(r, binary.BigEndian, v)
}

func (dec *Decoder) readNamedTag() (*NamedTag, error) {
	typ, err := dec.readType()
	if err != nil {
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
		err = dec.wrap(readBE(dec.r, &n))
		payload = n
	case TypeShort:
		var n int16
		err = dec.wrap(readBE(dec.r, &n))
		payload = n
	case TypeInt:
		var n int32
		err = dec.wrap(readBE(dec.r, &n))
		payload = n
	case TypeLong:
		var n int64
		err = dec.wrap(readBE(dec.r, &n))
		payload = n
	case TypeFloat:
		var x float32
		err = dec.wrap(readBE(dec.r, &x))
		payload = x
	case TypeDouble:
		var x float64
		err = dec.wrap(readBE(dec.r, &x))
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
		err = dec.errorf("unknown type (%v)", typ)
	}

	if err != nil {
		return nil, err
	}

	return &NamedTag{typ, name, payload}, nil
}

func (dec *Decoder) readType() (Type, error) {
	var typ Type
	err := dec.wrap(readBE(dec.r, &typ))
	return typ, err
}

func (dec *Decoder) readByteArray() ([]byte, error) {
	length, err := dec.readLength()
	if err != nil {
		return nil, err
	}

	b := make([]byte, length)
	if err := readBE(dec.r, b); err != nil {
		return nil, dec.wrap(err)
	}

	return b, nil
}

func (dec *Decoder) readLength() (int32, error) {
	var length int32
	err := dec.wrap(readBE(dec.r, &length))
	if length < 0 {
		err = dec.errorf("negative length (%d)", length)
	}
	return length, err
}

func (dec *Decoder) readString() (string, error) {
	var length int16
	if err := readBE(dec.r, &length); err != nil {
		return "", dec.wrap(err)
	}

	if length < 0 {
		return "", dec.errorf("negative length (%d)", length)
	}

	b := make([]byte, length)
	if err := readBE(dec.r, b); err != nil {
		return "", dec.wrap(err)
	}

	return string(b), nil
}

func (dec *Decoder) readList() (*List, error) {
	typ, err := dec.readType()
	if err != nil {
		return nil, err
	}

	length, err := dec.readLength()
	if err != nil {
		return nil, err
	}

	if typ == TypeEnd {
		return &List{}, nil
	}

	var array interface{}
	if typ < TypeByteArray {
		switch typ {
		case TypeByte:
			array = make([]int8, length)
		case TypeShort:
			array = make([]int16, length)
		case TypeInt:
			array = make([]int32, length)
		case TypeLong:
			array = make([]int64, length)
		case TypeFloat:
			array = make([]float32, length)
		case TypeDouble:
			array = make([]float64, length)
		}

		if err := readBE(dec.r, array); err != nil {
			return nil, dec.wrap(err)
		}

		return &List{typ, array}, nil
	}

	switch typ {
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
		return nil, dec.errorf("unknown type (%v)", typ)
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
			return nil, dec.errorf("duplicate name (%q)", tag.Name)
		}
		m[tag.Name] = &Tag{tag.Type, tag.Payload}
	}
}

func (dec *Decoder) readIntArray() ([]int32, error) {
	length, err := dec.readLength()
	if err != nil {
		return nil, err
	}

	a := make([]int32, length)
	if err := readBE(dec.r, a); err != nil {
		return nil, dec.wrap(err)
	}

	return a, nil
}

func (dec *Decoder) readLongArray() ([]int64, error) {
	length, err := dec.readLength()
	if err != nil {
		return nil, err
	}

	a := make([]int64, length)
	if err := readBE(dec.r, a); err != nil {
		return nil, dec.wrap(err)
	}

	return a, nil
}
