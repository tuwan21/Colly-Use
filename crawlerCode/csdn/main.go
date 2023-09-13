package main

import (
	"fmt"
	"github.com/gocolly/colly"
)

func main() {
	// 实例化默认收集器
	c := colly.NewCollector(
		colly.AllowedDomains("blog.csdn.net"),
	)
	c.OnHTML("article[class=\"blog-list-box\"] > a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		c.Visit(e.Request.AbsoluteURL(link))
	})

	// 收到response后调用
	c.OnResponse(func(r *colly.Response) {
		fmt.Printf("Response %s: %d bytes\n", r.Request.URL, len(r.Body))
	})

	// 在提出请求之前，请打印“正在访问...”
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// 请求过程中出现Error时调用
	c.OnError(func(c *colly.Response, err error) {
		fmt.Printf("Error %s: %v\n", c.Request.URL, err)
	})

	// 开始刮 https://hackerspaces.org
	c.Visit("https://blog.csdn.net/JN_HoweveR?spm=1000.2115.3001.5343")

}
