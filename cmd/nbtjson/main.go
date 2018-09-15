package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/njhanley/nbt"
	"github.com/pkg/errors"
)

func log(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format, v...)
}

func fatal(format string, v ...interface{}) {
	log(format, v...)
	os.Exit(1)
}

// NBT -> JSON
func decode(data []byte) ([]byte, error) {
	tag, err := nbt.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	b, err := json.Marshal(tag)
	return b, errors.WithStack(err)
}

// JSON -> NBT
func encode(data []byte) ([]byte, error) {
	var tag nbt.NamedTag
	if err := json.Unmarshal(data, &tag); err != nil {
		return nil, errors.WithStack(err)
	}
	buf := new(bytes.Buffer)
	err := nbt.Encode(buf, tag)
	return buf.Bytes(), errors.WithStack(err)
}

func gunzip(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer r.Close()

	b, err := ioutil.ReadAll(r)
	return b, errors.WithStack(err)
}

func main() {
	flag.Parse()

	if n := flag.NArg(); n < 2 {
		fatal("too few arguments\n")
	}

	const errfmt = "%q: %+v\n"

	switch flag.Arg(0) {
	case "decode", "dec", "d":
		filename := flag.Arg(1)

		data, err := ioutil.ReadFile(filename)
		if err != nil {
			fatal(errfmt, filename, err)
		}

		if b, err := gunzip(data); err != nil {
			if errors.Cause(err) != gzip.ErrHeader {
				fatal(errfmt, filename, err)
			}
		} else {
			data = b
		}

		text, err := decode(data)
		if err != nil {
			fatal(errfmt, filename, err)
		}

		if _, err := os.Stdout.Write(append(text, "\n"...)); err != nil {
			fatal(errfmt, filename, err)
		}
	case "encode", "enc", "e":
		filename := flag.Arg(1)

		data, err := ioutil.ReadFile(filename)
		if err != nil {
			fatal(errfmt, filename, err)
		}

		nbt, err := encode(data)
		if err != nil {
			fatal(errfmt, filename, err)
		}

		if _, err := os.Stdout.Write(nbt); err != nil {
			fatal(errfmt, filename, err)
		}
	default:
		fatal("unrecognized mode\n")
	}
}
