package nbt

import (
	"encoding/binary"
	"io"
	"math"
	"reflect"
	"runtime"

	"github.com/pkg/errors"
)

func Encode(w io.Writer, tag *NamedTag) (err error) {
	enc := &encoder{w: w}
	// handle possible panics from writeNamedTag and writeList
	defer func() {
		if v := recover(); v != nil {
			if e, ok := v.(*runtime.TypeAssertionError); ok {
				err = e
			} else {
				panic(v)
			}
		}
	}()
	return enc.writeNamedTag(tag)
}

type encoder struct {
	w io.Writer
}

func writeBE(w io.Writer, v interface{}) error {
	return binary.Write(w, binary.BigEndian, v)
}

func (enc *encoder) writeNamedTag(tag *NamedTag) (err error) {
	if err := writeBE(enc.w, tag.Type); err != nil {
		return err
	}

	if tag.Type == TypeEnd {
		return nil
	}

	if err := enc.writeString(tag.Name); err != nil {
		return err
	}

	switch tag.Type {
	case TypeByte, TypeShort, TypeInt, TypeLong, TypeFloat, TypeDouble:
		return writeBE(enc.w, tag.Payload)
	case TypeByteArray:
		return enc.writeByteArray(tag.Payload.([]byte))
	case TypeString:
		return enc.writeString(tag.Payload.(string))
	case TypeList:
		return enc.writeList(tag.Payload.(*List))
	case TypeCompound:
		return enc.writeCompound(tag.Payload.(Compound))
	case TypeIntArray:
		return enc.writeIntArray(tag.Payload.([]int32))
	case TypeLongArray:
		return enc.writeLongArray(tag.Payload.([]int64))
	default:
		return errors.Errorf("unknown type %#x", tag.Type)
	}
}

func (enc *encoder) writeByteArray(b []byte) error {
	length := len(b)
	if length > math.MaxInt32 {
		return errors.Errorf("length overflows int32 (%d)", length)
	}

	if err := writeBE(enc.w, int32(length)); err != nil {
		return err
	}

	return writeBE(enc.w, b)
}

func (enc *encoder) writeString(s string) error {
	length := len(s)
	if length > math.MaxInt16 {
		return errors.Errorf("length overflows int16 (%d)", length)
	}

	if err := writeBE(enc.w, int16(length)); err != nil {
		return err
	}

	return writeBE(enc.w, []byte(s))
}

func (enc *encoder) writeList(list *List) error {
	if err := writeBE(enc.w, list.Type); err != nil {
		return err
	}

	if list.Type == TypeEnd && list.Array == nil {
		return writeBE(enc.w, int32(0))
	}

	value := reflect.ValueOf(list.Array)
	if kind := value.Kind(); kind != reflect.Slice {
		return errors.Errorf("List.Array is not a slice (%v)", kind)
	}

	length := value.Len()
	if length > math.MaxInt32 {
		return errors.Errorf("length overflows int32 (%d)", length)
	}

	if err := writeBE(enc.w, int32(length)); err != nil {
		return err
	}

	switch list.Type {
	case TypeByte, TypeShort, TypeInt, TypeLong, TypeFloat, TypeDouble:
		return writeBE(enc.w, list.Array)
	case TypeByteArray:
		for _, a := range list.Array.([][]byte) {
			if err := enc.writeByteArray(a); err != nil {
				return err
			}
		}
	case TypeString:
		for _, a := range list.Array.([]string) {
			if err := enc.writeString(a); err != nil {
				return err
			}
		}
	case TypeList:
		for _, a := range list.Array.([]*List) {
			if err := enc.writeList(a); err != nil {
				return err
			}
		}
	case TypeCompound:
		for _, a := range list.Array.([]Compound) {
			if err := enc.writeCompound(a); err != nil {
				return err
			}
		}
	case TypeIntArray:
		for _, a := range list.Array.([][]int32) {
			if err := enc.writeIntArray(a); err != nil {
				return err
			}
		}
	case TypeLongArray:
		for _, a := range list.Array.([][]int64) {
			if err := enc.writeLongArray(a); err != nil {
				return err
			}
		}
	default:
		return errors.Errorf("unknown type %#x", list.Type)
	}

	return nil
}

func (enc *encoder) writeCompound(m Compound) error {
	for _, tag := range m {
		if err := enc.writeNamedTag(tag); err != nil {
			return err
		}
	}
	return enc.writeNamedTag(&NamedTag{})
}

func (enc *encoder) writeIntArray(a []int32) error {
	length := len(a)
	if length > math.MaxInt32 {
		return errors.Errorf("length overflows int32 (%d)", length)
	}

	if err := writeBE(enc.w, int32(length)); err != nil {
		return err
	}

	return writeBE(enc.w, a)
}

func (enc *encoder) writeLongArray(a []int64) error {
	length := len(a)
	if length > math.MaxInt32 {
		return errors.Errorf("length overflows int32 (%d)", length)
	}

	if err := writeBE(enc.w, int32(length)); err != nil {
		return err
	}

	return writeBE(enc.w, a)
}
