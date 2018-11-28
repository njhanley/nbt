package nbt

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
	"runtime"
	"sort"

	"github.com/pkg/errors"
)

type Encoder struct {
	w             io.Writer
	sortCompounds bool
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

func (enc *Encoder) Encode(tag *NamedTag) error {
	return enc.writeNamedTag(tag)
}

func (enc *Encoder) SortCompounds(on bool) {
	enc.sortCompounds = on
}

func (enc *Encoder) wrap(err error) error {
	return errors.WithStack(err)
}

func (enc *Encoder) errorf(format string, a ...interface{}) error {
	return enc.wrap(fmt.Errorf(format, a...))
}

func writeBE(w io.Writer, v interface{}) error {
	return binary.Write(w, binary.BigEndian, v)
}

func (enc *Encoder) writeNamedTag(tag *NamedTag) (err error) {
	// handle possible panics from type assertions in writeNamedTag and writeList
	defer func() {
		if v := recover(); v != nil {
			if e, ok := v.(*runtime.TypeAssertionError); ok {
				err = e
			} else {
				panic(v)
			}
		}
	}()

	if err := enc.writeType(tag.Type); err != nil {
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
		return enc.wrap(writeBE(enc.w, tag.Payload))
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
		return enc.errorf("unknown type (%v)", tag.Type)
	}
}

func (enc *Encoder) writeType(typ Type) error {
	return enc.wrap(writeBE(enc.w, typ))
}

func (enc *Encoder) writeByteArray(b []byte) error {
	if err := enc.writeLength(len(b)); err != nil {
		return err
	}
	return enc.wrap(writeBE(enc.w, b))
}

func (enc *Encoder) writeLength(length int) error {
	if length > math.MaxInt32 {
		return enc.errorf("length overflows int32 (%d)", length)
	}
	return enc.wrap(writeBE(enc.w, int32(length)))
}

func (enc *Encoder) writeString(s string) error {
	length := len(s)
	if length > math.MaxInt16 {
		return enc.errorf("length overflows int16 (%d)", length)
	}

	if err := writeBE(enc.w, int16(length)); err != nil {
		return enc.wrap(err)
	}

	return enc.wrap(writeBE(enc.w, []byte(s)))
}

func (enc *Encoder) writeList(list *List) error {
	if err := enc.writeType(list.Type); err != nil {
		return err
	}

	if list.Type == TypeEnd && list.Array == nil {
		return enc.writeLength(0)
	}

	value := reflect.ValueOf(list.Array)
	if kind := value.Kind(); kind != reflect.Slice {
		return enc.errorf("List.Array is not a slice (%v)", kind)
	}

	length := value.Len()
	if err := enc.writeLength(length); err != nil {
		return err
	}

	switch list.Type {
	case TypeByte, TypeShort, TypeInt, TypeLong, TypeFloat, TypeDouble:
		return enc.wrap(writeBE(enc.w, list.Array))
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
		return enc.errorf("unknown type (%v)", list.Type)
	}

	return nil
}

func (enc *Encoder) writeCompound(m Compound) error {
	if enc.sortCompounds {
		a := make([]*NamedTag, len(m))
		var i int
		for _, tag := range m {
			a[i] = tag
			i++
		}
		sort.Slice(a, func(i, j int) bool { return a[i].Name < a[j].Name })
		for _, tag := range a {
			if err := enc.writeNamedTag(tag); err != nil {
				return err
			}
		}
	} else {
		for _, tag := range m {
			if err := enc.writeNamedTag(tag); err != nil {
				return err
			}
		}
	}
	return enc.writeNamedTag(&NamedTag{})
}

func (enc *Encoder) writeIntArray(a []int32) error {
	if err := enc.writeLength(len(a)); err != nil {
		return err
	}
	return enc.wrap(writeBE(enc.w, a))
}

func (enc *Encoder) writeLongArray(a []int64) error {
	if err := enc.writeLength(len(a)); err != nil {
		return err
	}
	return enc.wrap(writeBE(enc.w, a))
}
