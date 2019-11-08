package main

import (
	"flag"
	"fmt"
	"hammer-x/config"
	"hammer-x/download/youtube"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

func init() {
	flag.StringVar(&config.Proxy, "x", "", "use proxy to flip over gfw")
	flag.StringVar(&config.Socket, "s", "", "use socket to flip over gfw")
	flag.StringVar(&config.Uri, "i", "", "download url of video ")
	flag.BoolVar(&config.Help, "h", false, "this help")
	flag.BoolVar(&config.Debug, "d", false, "debug mode")
	flag.Usage = usage
}

const TOTAL = 20

func main() {

	var progress = 0
	var position = 1
Loop:
	for {
		if progress > 0 {
			fmt.Printf("\033[%dA\033[K", position)
		}

		output := fmt.Sprintf(
			"%s%s%s",
			"progress: ",
			strings.Repeat("=", progress),
			strings.Repeat("-", TOTAL-progress),
		)

		fmt.Printf("%s \033[K\n", output)
		//fmt.Printf("\033[%dA\033[K", 1)
		if progress >= 20 {
			break Loop
		}
		progress++
		time.Sleep(time.Duration(200) * time.Millisecond)
	}
	return
	flag.Parse()
	download_url := config.Uri
	usr, _ := user.Current()
	dir := fmt.Sprintf("%v/Movies/hammer-x", usr.HomeDir)
	log.Println("download to dir=", dir)

	if strings.Contains(download_url, "youtube") {
		y := youtube.InitYoutube(config.Debug)
		y.Agency = config.Proxy
		if err := y.DecodeURL(download_url); err != nil {
			fmt.Println("err:", err)
		}
		if err := y.StartDownload(filepath.Join(dir, y.VideoID+".mp4")); err != nil {
			fmt.Println("err:", err)
		}
	}
	if config.Help {
		flag.Usage()
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `
Usage: hammer-x  [-x proxy] [-i url]
Options:
`)
	flag.PrintDefaults()
}
