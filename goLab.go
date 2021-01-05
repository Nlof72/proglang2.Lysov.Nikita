package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

// WriteCounter Структура для райтера
type WriteCounter struct {
	Total uint64
	Time  uint64
	Speed uint64
}

// Вызывается в TeeReader
func (wc *WriteCounter) Write(p []byte) (int, error) {

	n := len(p)
	wc.Total += uint64(n)
	return n, nil
}

// Печать прогресса
func progress(wc *WriteCounter, leng uint64) {
	for {
		fmt.Println("Download:", wc.Total, "bytes", "with middle speed:", wc.Speed, "bytes/s")
		time.Sleep(time.Second)
		wc.Time++
		wc.Speed = wc.Total / wc.Time
		if wc.Total == leng {
			fmt.Println(wc.Total)
			return
		}
	}
}

func main() {
	// должно задействовать 4 ядра
	runtime.GOMAXPROCS(4)

	fileURL := ""
	fmt.Println("Enter URL adress for downloading")
	fmt.Scan(&fileURL)
	fmt.Println("Download Started")
	err := DownloadFile(fileURL[strings.LastIndex(fileURL, "/")+1:len(fileURL)], fileURL)
	if err != nil {
		panic(err)
	}

	fmt.Println("Download Finished")
}

// DownloadFile функция загрузки
func DownloadFile(filepath string, url string) error {
	today := time.Now()

	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Ошибка в ссылке")
		out.Close()
		return err
	}
	defer resp.Body.Close()

	leng := uint64(resp.ContentLength)
	counter := &WriteCounter{}
	go progress(counter, leng)

	readerr := io.TeeReader(resp.Body, counter)

	if _, err = io.Copy(out, readerr); err != nil {
		out.Close()
		return err
	}

	today1 := time.Now()
	dowtime := today1.Sub(today)
	fmt.Println("Download file:", filepath, "with size:", counter.Total, "bytes", "for:", dowtime, "with middle speed:", counter.Speed, "bytes/s")

	fmt.Print("\n")

	out.Close()

	if err = os.Rename(filepath+".tmp", filepath); err != nil {
		return err
	}
	return nil
}
