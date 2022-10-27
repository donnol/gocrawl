package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// curl
// -x http://192.168.232.1:7890
// -H 'user-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36'
// -v
// https://mmzztt.com

func main() {
	url := "https://mmzztt.com/"
	images := crawl(url)
	fmt.Printf("images: %+v\n", images)
	// http.Transport.pro
}

var (
	httpClient = func() *http.Client {

		tp := http.DefaultTransport.(*http.Transport)
		tp.Proxy = func(r *http.Request) (*url.URL, error) {
			link, err := url.Parse("http://192.168.232.1:7890")
			if err != nil {
				return nil, err
			}
			return link, nil
		}

		return &http.Client{
			Transport: tp,
			Timeout:   10 * time.Second,
		}
	}()
)

func crawl(url string) []string {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36")

	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	images := []string{}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		singleImage, ok := s.Find("img").Attr("data-src")
		if !ok {
			return
		}
		images = append(images, singleImage)

		// 保存到文件
		saveImage(singleImage)
	})
	if len(images) == 0 {
		html, err := doc.Html()
		if err != nil {
			panic(err)
		}
		fmt.Printf("doc: %s\n", html)
	}

	return images
}

// 将图片下载并保存到本地
func saveImage(url string) {
	// 图片内容
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36")
	res, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

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
	defer dst.Close()

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
