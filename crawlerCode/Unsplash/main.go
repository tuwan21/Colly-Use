package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"sync/atomic"
	"time"
)

func main() {
	c1 := colly.NewCollector(
		colly.Async(true),
	)
	c2 := c1.Clone()
	c3 := c1.Clone()

	/*	c1.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})*/
	err := c1.Limit(&colly.LimitRule{
		DomainRegexp: `unsplash\.com`,
		RandomDelay:  500 * time.Millisecond,
		Parallelism:  12,
	})
	if err != nil {
		return
	}

	c1.OnHTML("div.zmDAx a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if link == "" {
			return
		}
		c2.Visit(e.Request.AbsoluteURL(link))
	})

	c2.OnHTML("div.MorZF > img[srcset]", func(e *colly.HTMLElement) {
		src := e.Attr("srcset")
		if src == "" {
			return
		}
		c3.Visit(src)
	})
	var count uint32
	c3.OnResponse(func(r *colly.Response) {
		fileName := fmt.Sprintf("images/img%d.jpg", atomic.AddUint32(&count, 1))
		err := r.Save(fileName)
		if err != nil {
			fmt.Printf("saving %s failed:%v\n", fileName, err)
		} else {
			fmt.Printf("saving %s success\n", fileName)
		}
	})

	c1.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c1.OnError(func(r *colly.Response, err error) {
		fmt.Println("Visiting", r.Request.URL, "failed:", err)
	})
	c3.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c1.Visit("https://unsplash.com/backgrounds")
	c1.Wait()
	c2.Wait()
	c3.Wait()
}
