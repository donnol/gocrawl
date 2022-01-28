package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	url := "https://mmzztt.com/"
	crawl(url)
}

func crawl(url string) []string {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		panic(err)
	}

	images := []string{}
	doc.Find(".postlist li").Each(func(i int, s *goquery.Selection) {
		singleImage, ok := s.Find("img").Attr("data-original")
		if !ok {
			return
		}
		images = append(images, singleImage)

		// 保存到文件
		saveImage(singleImage)
	})

	return images
}

// 将图片下载并保存到本地
func saveImage(url string) {
	// 图片内容
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	if err != nil {
		panic(err)
	}

	// 目录
	Dirname := "./tmp/" + time.Now().Format("2006-01-02") + "/"
	if !isDirExist(Dirname) {
		err = os.MkdirAll(Dirname, 0755)
		if err != nil {
			panic(err)
		}
		fmt.Printf("dir %s created\n", Dirname)
	}

	// 文件
	filename := filepath.Base(url)
	dst, err := os.Create(Dirname + filename)
	if err != nil {
		panic(err)
	}

	// 写入文件
	io.Copy(dst, res.Body)
}

func isDirExist(path string) bool {
	p, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		return p.IsDir()
	}
}
