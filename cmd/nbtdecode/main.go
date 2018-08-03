package main

import (
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/njhanley/nbt"
)

func log(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format, v...)
}

func fatal(format string, v ...interface{}) {
	log(format, v...)
	os.Exit(1)
}

func isGZIP(f *os.File) (bool, error) {
	sig := make([]byte, 2)

	if _, err := io.ReadFull(f, sig); err != nil {
		return false, err
	}

	if _, err := f.Seek(0, 0); err != nil {
		return false, err
	}

	return sig[0] == 0x1f && sig[1] == 0x8b, nil
}

func decodeFile(name string) (nbt.NamedTag, error) {
	file, err := os.Open(name)
	if err != nil {
		return nbt.NamedTag{}, err
	}
	defer file.Close()

	gzipped, err := isGZIP(file)
	if err != nil {
		return nbt.NamedTag{}, err
	}

	var r io.Reader
	if gzipped {
		gz, err := gzip.NewReader(file)
		if err != nil {
			return nbt.NamedTag{}, err
		}
		defer gz.Close()

		r = gz
	} else {
		r = file
	}

	return nbt.Decode(r)
}

type formatter func(nbt.NamedTag) error

func goSyntaxFormat(tag nbt.NamedTag) error {
	fmt.Printf("%#v\n", tag)
	return nil
}

func jsonFormat(tag nbt.NamedTag) error {
	b, err := json.Marshal(tag)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", b)
	return nil
}

func jsonIndentFormat(tag nbt.NamedTag) error {
	b, err := json.MarshalIndent(tag, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", b)
	return nil
}

var (
	indent   = flag.Bool("i", false, "indent JSON")
	goSyntax = flag.Bool("g", false, "output Go syntax representation instead of JSON")
)

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		fatal("no file(s) specified\n")
	}

	var output formatter
	if *goSyntax {
		output = goSyntaxFormat
	} else {
		if *indent {
			output = jsonIndentFormat
		} else {
			output = jsonFormat
		}
	}

	for _, arg := range flag.Args() {
		tag, err := decodeFile(arg)
		if err != nil {
			log("%q: %v\n", arg, err)
			continue
		}

		if err = output(tag); err != nil {
			log("%q: %v\n", arg, err)
			continue
		}
	}
}
