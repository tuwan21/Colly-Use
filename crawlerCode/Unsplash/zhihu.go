package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"time"
)

type Jianshu struct {
	B string `selector:"div._2mYfmT"`
	C string `selector:"a._1qp91i _1OhGeD"`
	D string `selector:"img._13D2Eh"`
}

func main() {
	c1 := colly.NewCollector(
		colly.Async(true),
	)
	err := c1.Limit(&colly.LimitRule{
		DomainRegexp: `jianshu\.com`,
		RandomDelay:  500 * time.Millisecond,
		Parallelism:  12,
	})
	if err != nil {
		return
	}

	var jianshus []Jianshu
	c1.OnHTML("div.rEsl9f", func(e *colly.HTMLElement) {
		jianshu := &Jianshu{}
		if err := e.Unmarshal(jianshu); err != nil {
			fmt.Println("error:", err)
			return
		}
		jianshus = append(jianshus, *jianshu)
		attr := e.Attr("href")
		fmt.Printf("Link found: %q -> %s\n", e.Text, attr)
	})

	c1.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c1.OnError(func(r *colly.Response, err error) {
		fmt.Println("Visiting", r.Request.URL, "failed:", err)
	})
	c1.Visit("https://www.jianshu.com/p/0d08d1d09319")
	c1.Wait()
}
