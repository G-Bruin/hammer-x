package utils

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

type PassThru struct {
	io.Reader
	total  int64
	level  int
	Length int64
	full   int
}

func (pt *PassThru) Read(p []byte) (int, error) {
	pt.full = 100
	n, err := pt.Reader.Read(p)
	pt.total += int64(n)
	if err == nil {
		fmt.Println("Read", pt.total, "bytes for a total of", pt.Length)
	}

	strInt64 := strconv.FormatInt((pt.total * 100 / pt.Length), 10)
	percent, _ := strconv.Atoi(strInt64)

	if percent > 0 {
		fmt.Printf("\033[%dA\033[K", percent)
	}

	output := fmt.Sprintf(
		"%s%s",
		"progress: "+strInt64+"% ",
		strings.Repeat("=", percent),
	)

	fmt.Printf("%s \033[K\n", output)

	//fmt.Printf("%v\n",result)
	//a := pt.total/pt.Length
	//percent := ((pt.total / pt.Length) * 100)
	//if (y.downloadLevel <= currentPercent) && (y.downloadLevel < 100) {
	//	y.downloadLevel++
	//	y.DownloadPercent <- int64(y.downloadLevel)
	//}
	return n, err
}
