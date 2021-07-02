package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const URL string = "https://en.wikipedia.org/wiki/Wikipedia:Picture_of_the_day"

func updateReadme(link string) {
	// read common
	common, err := os.ReadFile("common.md")
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("README.md")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if _, err := f.Write(common); err != nil {
		log.Fatal(err)
	}

	linkMd := fmt.Sprintf("\n![](%s)", link)
	if _, err := f.WriteString(linkMd); err != nil {
		log.Fatal(err)
	}
}

func getPOTD() string {
	var link string
	// Request the HTML page.
	res, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("BUG: status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("a[class=image]>img").EachWithBreak(func(i int, s *goquery.Selection) bool {
		srcset, ok := s.Attr("srcset")
		if ok {
			imageURLs := strings.Split(srcset, " ")
			link = "https:" + imageURLs[len(imageURLs)-2]
			fmt.Printf("link of POTD: %s", link)
			return false
		}
		return true
	})
	if link == "" {
		log.Fatal("BUG: link is empty")
	}
	return link
}

func main() {
	link := getPOTD()
	updateReadme(link)
}
