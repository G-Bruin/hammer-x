package main

import (
	"bytes"
	"flag"
	"fmt"
	progressbar "github.com/schollz/progressbar/v2"
	"hammer-x/config"
	"hammer-x/download/youtube"
	"hammer-x/utils"
	"io"
	"log"
	"net/http"
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

func main() {

	fmt.Println("downloading go1.12.5.linux-amd64.tar.gz")
	defer os.Remove("go1.12.5.linux-amd64.tar.gz")
	urlToGet := "https://dl.google.com/go/go1.12.5.linux-amd64.tar.gz"
	req, _ := http.NewRequest("GET", urlToGet, nil)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	var out io.Writer
	f, _ := os.OpenFile("go1.12.5.linux-amd64.tar.gz", os.O_CREATE|os.O_WRONLY, 0644)
	out = f
	defer f.Close()

	bar := progressbar.NewOptions(
		int(resp.ContentLength),
		progressbar.OptionSetBytes(int(resp.ContentLength)),
		progressbar.OptionThrottle(10*time.Millisecond),
	)
	out = io.MultiWriter(out, bar)
	io.Copy(out, resp.Body)
	fmt.Println("done")
	//bar := progressbar.NewOptions(
	//	int(resp.ContentLength),
	//	progressbar.OptionSetBytes(int(resp.ContentLength)),
	//)
	//out = io.MultiWriter(out, bar)
	//io.Copy(out, resp.Body)
	//	var progress = 0
	//	var position = 1
	//Loop:
	//	for {
	//		if progress > 0 {
	//			fmt.Printf("\033[%dA\033[K", position)
	//		}
	//
	//		output := fmt.Sprintf(
	//			"%s%s%s",
	//			"progress: ",
	//			strings.Repeat("=", progress),
	//			strings.Repeat("-", TOTAL-progress),
	//		)
	//
	//		fmt.Printf("%s \033[K\n", output)
	//		//fmt.Printf("\033[%dA\033[K", 1)
	//		if progress >= 20 {
	//			break Loop
	//		}
	//		progress++
	//		time.Sleep(time.Duration(200) * time.Millisecond)
	//	}

	var src io.Reader    // Source file/url/etc
	var dst bytes.Buffer // Destination file/buffer/etc

	src = bytes.NewBufferString(strings.Repeat("Some random input data", 1000))

	src = &utils.PassThru{Reader: src, Length: 22000}

	fmt.Println("sadadsad")
	count, err := io.Copy(&dst, src)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Transferred", count, "bytes")

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
