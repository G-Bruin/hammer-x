package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"v-download/utils"
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
	download_url := u

	usr, _ := user.Current()
	dir := fmt.Sprintf("%v/Movies/v-download", usr.HomeDir)
	log.Println("download to dir=", dir)

	if strings.Contains(download_url, "youtube") {
		youtube := utils.InitYoutube(true)
		youtube.Agency = p
		if err := youtube.DecodeURL(download_url); err != nil {
			fmt.Println("err:", err)
		}
		if err := youtube.StartDownload(filepath.Join(dir, youtube.VideoID+".mp4")); err != nil {
			fmt.Println("err:", err)
		}
	}

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
