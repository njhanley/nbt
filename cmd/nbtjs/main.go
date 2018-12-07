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

var options struct {
	indent        string
	revert        bool
	sortCompounds bool
	gzip          bool
	gzipLevel     int
	verbose       bool
}

func init() {
	flag.StringVar(&options.indent, "i", "", "indent output JSON with string")
	flag.BoolVar(&options.revert, "r", false, "revert JSON to NBT")
	flag.BoolVar(&options.sortCompounds, "s", false, "write compound tags in lexically sorted order")
	flag.BoolVar(&options.gzip, "z", false, "gzip the output NBT")
	flag.IntVar(&options.gzipLevel, "zlevel", 6, "gzip compression level, 0 = none, 1 = fast, 9 = best")
	flag.BoolVar(&options.verbose, "v", false, "verbose mode")
}

type exitCode int

func exit(code int) {
	panic(exitCode(code))
}

func handleExit() {
	if v := recover(); v != nil {
		if code, ok := v.(exitCode); ok {
			os.Exit(int(code))
		}
		panic(v)
	}
}

func info(prefix string, v interface{}) {
	if e, ok := v.(*os.PathError); ok {
		v = e.Err
	}

	if options.verbose {
		fmt.Fprintf(os.Stderr, "%s: %+v\n", prefix, v)
	} else {
		fmt.Fprintf(os.Stderr, "%s: %v\n", prefix, v)
	}
}

func fatal(prefix string, v interface{}) {
	info(prefix, v)
	exit(1)
}

func closeIO(c io.Closer, name string) {
	if err := c.Close(); err != nil {
		fatal(name, err)
	}
}

func nbtToJSON(in *os.File, out *os.File) {
	var dec *nbt.Decoder
	if r, err := gzip.NewReader(in); err != gzip.ErrHeader {
		if options.verbose {
			info(in.Name(), "decompressing")
		}

		if err != nil {
			fatal(in.Name(), err)
		}
		defer closeIO(r, in.Name())

		dec = nbt.NewDecoder(r)
	} else {
		dec = nbt.NewDecoder(in)
	}

	tag, err := dec.Decode()
	if err != nil {
		fatal(in.Name(), err)
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", options.indent)

	if err := enc.Encode(tag); err != nil {
		fatal(out.Name(), err)
	}
}

func jsonToNBT(in *os.File, out *os.File) {
	dec := json.NewDecoder(in)

	tag := new(nbt.NamedTag)
	if err := dec.Decode(tag); err != nil {
		fatal(in.Name(), err)
	}

	var enc *nbt.Encoder
	if options.gzip {
		w, err := gzip.NewWriterLevel(out, options.gzipLevel)
		if err != nil {
			fatal(out.Name(), err)
		}
		defer closeIO(w, out.Name())

		enc = nbt.NewEncoder(w)
	} else {
		enc = nbt.NewEncoder(out)
	}

	enc.SortCompounds(options.sortCompounds)

	if err := enc.Encode(tag); err != nil {
		fatal(out.Name(), err)
	}
}

func main() {
	defer handleExit()

	flag.Parse()

	var infile, outfile string
	switch flag.NArg() {
	default:
		flag.Usage()
		exit(2)
	case 2:
		if arg := flag.Arg(1); arg != "-" {
			outfile = arg
		}
		fallthrough
	case 1:
		if arg := flag.Arg(0); arg != "-" {
			infile = arg
		}
	case 0:
	}

	var in *os.File
	if infile == "" {
		in = os.Stdin
	} else {
		file, err := os.Open(infile)
		if err != nil {
			fatal(infile, err)
		}
		defer closeIO(file, infile)

		in = file
	}

	var out *os.File
	if outfile == "" {
		out = os.Stdout
	} else {
		file, err := os.Create(outfile)
		if err != nil {
			fatal(outfile, err)
		}
		defer closeIO(file, outfile)

		out = file
	}

	if options.revert {
		jsonToNBT(in, out)
	} else {
		nbtToJSON(in, out)
	}
}
