package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type Page struct {
	PageNum  int `form:"page_num"`
	PageSize int `form:"page_size"`
	Keyword  int `form:"keyword"`
	Desc     int `form:"desc"`
}

func main() {
	c := gin.Context{}
	var p Page
	if err := c.ShouldBindQuery(&p); err != nil {
		fmt.Println("参数错误")
		return
	}
	if p.PageNum <= 0 {
		p.PageNum = 1
	}

}
