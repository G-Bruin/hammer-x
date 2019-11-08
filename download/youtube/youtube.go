package youtube

import (
	"errors"
	"fmt"
	"hammer-x/utils"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func SetLogOutput(w io.Writer) {
	log.SetOutput(w)
}

func InitYoutube(debug bool) *Youtube {
	return &Youtube{DebugMode: debug, DownloadPercent: make(chan int64, 100)}
}

type stream map[string]string

type Youtube struct {
	DebugMode         bool
	StreamList        []stream
	VideoID           string
	Agency            string
	videoInfo         string
	DownloadPercent   chan int64
	contentLength     int64
	totalWrittenBytes float64
	downloadLevel     float64
}

func (y *Youtube) DecodeURL(url string) error {
	err := y.findVideoID(url)
	if err != nil {
		return fmt.Errorf("findVideoID error=%s", err)
	}

	err = y.getVideoInfo()
	if err != nil {
		return fmt.Errorf("getVideoInfo error=%s", err)
	}
	err = y.parseVideoInfo()
	if err != nil {
		return fmt.Errorf("parse video info failed, err=%s", err)
	}

	return nil
}

func (y *Youtube) StartDownload(destFile string) error {
	//download highest resolution on [0]
	var err error
	for _, v := range y.StreamList {
		url := v["url"]
		y.log(fmt.Sprintln("Download url=", url))
		y.log(fmt.Sprintln("Download to file=", destFile))
		err = y.videoDLWorker(destFile, url)
		if err == nil {
			break
		}
	}
	fmt.Println("download finished")
	return err
}

func (y *Youtube) parseVideoInfo() error {
	answer, err := url.ParseQuery(y.videoInfo)
	if err != nil {
		return err
	}
	status, ok := answer["status"]
	if !ok {
		err = fmt.Errorf("no response status found in the server's answer")
		return err
	}
	if status[0] == "fail" {
		reason, ok := answer["reason"]
		if ok {
			err = fmt.Errorf("'fail' response status found in the server's answer, reason: '%s'", reason[0])
		} else {
			err = errors.New(fmt.Sprint("'fail' response status found in the server's answer, no reason given"))
		}
		return err
	}
	if status[0] != "ok" {
		err = fmt.Errorf("non-success response status found in the server's answer (status: '%s')", status)
		return err
	}

	// read the streams map
	streamMap, ok := answer["url_encoded_fmt_stream_map"]
	if !ok {
		err = errors.New(fmt.Sprint("no stream map found in the server's answer"))
		return err
	}
	// read each stream
	streamsList := strings.Split(streamMap[0], ",")
	var streams []stream
	for streamPos, streamRaw := range streamsList {
		streamQry, err := url.ParseQuery(streamRaw)
		if err != nil {
			log.Printf("An error occured while decoding one of the video's stream's information: stream %d: %s\n", streamPos, err)
			continue
		}
		streams = append(streams, stream{
			"quality": streamQry["quality"][0],
			"type":    streamQry["type"][0],
			"url":     streamQry["url"][0],
		})
		y.log(fmt.Sprintf("Stream found: quality '%s', format '%s'", streamQry["quality"][0], streamQry["type"][0]))
	}
	y.StreamList = streams
	return nil
}

func (y *Youtube) getVideoInfo() error {
	target_url := "http://youtube.com/get_video_info?video_id=" + y.VideoID
	body, _ := utils.Get(target_url, "", nil)
	y.videoInfo = body
	return nil
}

func (y *Youtube) findVideoID(url string) error {
	videoID := url
	if strings.Contains(videoID, "youtu") || strings.ContainsAny(videoID, "\"?&/<%=") {
		reList := []*regexp.Regexp{
			regexp.MustCompile(`(?:v|embed|watch\?v)(?:=|/)([^"&?/=%]{11})`),
			regexp.MustCompile(`(?:=|/)([^"&?/=%]{11})`),
			regexp.MustCompile(`([^"&?/=%]{11})`),
		}
		for _, re := range reList {
			if isMatch := re.MatchString(videoID); isMatch {
				subs := re.FindStringSubmatch(videoID)
				videoID = subs[1]
			}
		}
	}
	log.Printf("Found video id: '%s'", videoID)
	y.VideoID = videoID
	if strings.ContainsAny(videoID, "?&/<%=") {
		return errors.New("invalid characters in video id")
	}
	if len(videoID) < 10 {
		return errors.New("the video id must be at least 10 characters long")
	}
	return nil
}

func (y *Youtube) Write(p []byte) (n int, err error) {
	n = len(p)
	//y.totalWrittenBytes = y.totalWrittenBytes + float64(n)
	//currentPercent := ((y.totalWrittenBytes / y.contentLength) * 100)
	//if (y.downloadLevel <= currentPercent) && (y.downloadLevel < 100) {
	//	y.downloadLevel++
	//	y.DownloadPercent <- int64(y.downloadLevel)
	//}

	return
}

func (y *Youtube) videoDLWorker(destFile string, target string) error {
	resp, err := utils.Request("GET", target, nil, nil)
	defer resp.Body.Close()
	y.contentLength = resp.ContentLength
	if resp.StatusCode != 200 {
		return errors.New("non 200 status code received")
	}
	err = os.MkdirAll(filepath.Dir(destFile), 0755)
	if err != nil {
		return err
	}
	out, err := os.Create(destFile)
	if err != nil {
		return err
	}
	mw := io.MultiWriter(out, y)
	src := &utils.PassThru{Reader: resp.Body, Length: y.contentLength}
	_, err = io.Copy(mw, src)
	if err != nil {
		log.Println("download video err=", err)
		return err
	}
	return nil
}

func (y *Youtube) log(logText string) {
	if y.DebugMode {
		log.Println(logText)
	}
}
