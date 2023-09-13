package main

import (
	"fmt"
	"github.com/gocolly/colly"
)

type Csdn struct {
	Name string `json:"name"`
}

func main() {
	csdn := Csdn{}
	c1 := colly.NewCollector(
		colly.Async(true),
	)
	c2 := c1.Clone()
	c1.OnHTML("article.blog-list-box a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		c2.Visit(e.Request.AbsoluteURL(link))
	})

	c2.OnHTML("div.article-title-box h1", func(e *colly.HTMLElement) {
		text := e.Text
		fmt.Println("text", text)
		csdn.Name = text
	})

	c1.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	c1.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error %s: %v\n", r.Request.URL, err)
	})
	c1.Visit("https://blog.csdn.net/JN_HoweveR?spm=1000.2115.3001.5343")
	c1.Wait()
	c2.Wait()
}
