package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	h bool
	p string
	u string
)

func init() {
	flag.BoolVar(&h, "h", false, "this help")
	flag.StringVar(&p, "p", "", "use proxy to flip over gfw")
	flag.StringVar(&u, "u", "", "download url of video ")

	flag.Usage = usage
}

func main() {
	flag.Parse()
	fmt.Println("proxy:", p)
	fmt.Println("url:", u)
	if h {
		flag.Usage()
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `
Usage: v-download  [-p proxy] [-u url]
Options:
`)
	flag.PrintDefaults()
}