package main

import (
	"net/http"
	"io"
	"os"
	"strings"
	"time"
	"fmt"
)

type Progress struct {
	BytesCount int
	IsStopped  bool
}

func NewProgress(reader io.Reader) (io.Reader, *Progress) {
	progress := Progress{
		BytesCount: 0,
		IsStopped:  false,
	}
	teeReader := io.TeeReader(reader, &progress)
	return teeReader, &progress
}

func (progress *Progress) Write(bytes []byte) (n int, err error) {
	progress.BytesCount += len(bytes)
	return len(bytes), nil
}

func (progress *Progress) StartTick() {
	go func ()  {
		for {
			time.Sleep(time.Second)

			if progress.IsStopped{
				break
			}

			fmt.Println("Downloaded", progress.BytesCount/1024/1024, "MB")
		}
	}()
}

func (progress *Progress) StopTick() {
	progress.IsStopped = true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No URL")
		os.Exit(1)
	}

	URL := os.Args[1]

	if resp, err := http.Get(URL); err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	} else if resp.StatusCode != 200 {
		fmt.Println(resp.Status)
		os.Exit(3)
	} else {
		parts := strings.Split(URL, "/")
		name := strings.Split(parts[len(parts) - 1], "?")[0]

		if file, err := os.Create(name); err != nil {
			fmt.Println(err.Error())
			os.Exit(4)
		} else {
			defer file.Close()

			reader, progress := NewProgress(resp.Body)

			progress.StartTick()

			if size, err := io.Copy(file, reader); err != nil {
				fmt.Println(err.Error())
				os.Exit(5)
			} else {
				progress.StopTick()
				fmt.Println("File downloading is done. Size:", size/1024/1024, "MB")
			}
		}
	}
}
