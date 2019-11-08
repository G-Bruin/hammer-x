package utils

import (
	"fmt"
	"io"
)

type PassThru struct {
	io.Reader
	total  int64
	level  int
	Length int64
}

func (pt *PassThru) Read(p []byte) (int, error) {
	n, err := pt.Reader.Read(p)
	pt.total += int64(n)
	if err == nil {
		fmt.Println("Read", pt.total, "bytes for a total of", pt.Length)
	}
	//currentPercent := ((pt.total / pt.Length) * 100)
	//if (y.downloadLevel <= currentPercent) && (y.downloadLevel < 100) {
	//	y.downloadLevel++
	//	y.DownloadPercent <- int64(y.downloadLevel)
	//}
	return n, err
}
