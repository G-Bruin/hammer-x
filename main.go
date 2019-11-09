package main

import (
	"flag"
	"fmt"
	"hammer-x/config"
	"hammer-x/download/youtube"
	"hammer-x/utils"
	"os"
	"os/user"
	"path/filepath"
)

func init() {
	flag.StringVar(&config.Proxy, "x", "", "use proxy to flip over gfw")
	flag.StringVar(&config.Socket, "s", "", "use socket to flip over gfw")
	flag.StringVar(&config.Uri, "i", "", "download url of video ")
	flag.BoolVar(&config.Help, "h", false, "this help")
	flag.BoolVar(&config.Debug, "d", false, "debug mode")
	flag.Usage = usage
}

func main() {
	flag.Parse()
	if config.Help {
		flag.Usage()
		return
	}

	usr, _ := user.Current()
	dir := fmt.Sprintf("%v/Movies/hammer-x", usr.HomeDir)
	fmt.Println("download to dir is ", dir)
	host := utils.FindDomain(config.Uri)

	switch host {
	case "youtube":
		y := youtube.InitYoutube(config.Debug)
		y.Agency = config.Proxy
		if err := y.DecodeURL(config.Uri); err != nil {
			fmt.Println("err:", err)
		}
		if err := y.StartDownload(filepath.Join(dir, y.VideoID+".mp4")); err != nil {
			fmt.Println("err:", err)
		}
	default:
		fmt.Println("请输入url")
		return
	}
}

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, `
Usage: hammer-x  [-x proxy] [-i url]
Options:
`)
	flag.PrintDefaults()
}
