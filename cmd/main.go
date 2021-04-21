package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"

	"github.com/wailovet/impl"
)

const usage = `impl [-dir directory] <recv> <iface>

impl generates method stubs for recv to implement iface.

Examples:

impl 'f *File' io.Reader
impl Murmur hash.Hash
impl -dir $GOPATH/src/github.com/josharian/impl Murmur hash.Hash

Don't forget the single quotes around the receiver type
to prevent shell globbing.
`

var (
	flagSrcDir = flag.String("dir", "", "package source directory, useful for vendored code")
)

// validReceiver reports whether recv is a valid receiver expression.
func validReceiver(recv string) bool {
	if recv == "" {
		// The parse will parse empty receivers, but we don't want to accept them,
		// since it won't generate a usable code snippet.
		return false
	}
	fset := token.NewFileSet()
	_, err := parser.ParseFile(fset, "", "package hack\nfunc ("+recv+") Foo()", 0)
	return err == nil
}
func main() {
	flag.Parse()

	if len(flag.Args()) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}

	recv, iface := flag.Arg(0), flag.Arg(1)
	if !validReceiver(recv) {
		fatal(fmt.Sprintf("invalid receiver: %q", recv))
	}

	src, err := impl.Do(recv, iface, *flagSrcDir)
	if err != nil {
		log.Panic(err)
	}
	fmt.Print(string(src))
}

func fatal(msg interface{}) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
