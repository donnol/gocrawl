package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	link = "https://mmzztt.com/"
)

var (
	chars = []string{
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s",
		"t", "u", "v", "w", "x", "y", "z",
	}
)

func main() {
	// crawl(link)
	// saveImage("https://p.iimzt.com/2022/01/20s05h.jpg")
	// https://p.iimzt.com/2022/01/20s02i.jpg

	for i := 20; i < 30; i++ {
		for j := 1; j < 20; j++ {
			// for _, cb := range chars {
			// 	_ = cb
			for _, cf := range chars {

				name := fmt.Sprintf("%02d%s%02d%s", i, "s", j, cf)
				url := "https://p.iimzt.com/2022/01/" + name + ".jpg"
				saveImage(url)

				time.Sleep(250 * time.Millisecond)
			}
			// }
		}
	}
}

func crawl(url string) []string {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		panic(err)
	}

	images := []string{}
	doc.Find(".uk-grid li").Each(func(i int, s *goquery.Selection) {
		singleImage, ok := s.Find("img").Attr("data-src")
		fmt.Printf("image: %s\n", singleImage)
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
	client := &http.Client{
		Transport: http.DefaultTransport,
		// CheckRedirect: func(*http.Request, []*http.Request) error { panic("not implemented") },
		// Jar:           nil,
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("new request failed: %+v\n", err)
		return
	}
	req.Header.Add("Referer", link)
	// 图片内容
	res, err := client.Do(req)
	if err != nil {
		log.Printf("client do failed: %+v\n", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode == 404 || res.StatusCode == 429 {
		fmt.Printf("http status code: %d, url: %s\n", res.StatusCode, url)
		return
	}

	// 目录
	Dirname := "./tmp/" + time.Now().Format("2006-01-02") + "/"
	if !isDirExist(Dirname) {
		err = os.MkdirAll(Dirname, 0755)
		if err != nil {
			log.Printf("mkdir all failed: %+v\n", err)
			return
		}
		fmt.Printf("dir %s created\n", Dirname)
	}

	// 文件
	filename := filepath.Base(url)
	ext := filepath.Ext(filename)
	n := rand.Int()
	filename = strconv.Itoa(n) + ext
	dst, err := os.Create(Dirname + filename)
	if err != nil {
		log.Printf("create file failed: %+v\n", err)
		return
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
