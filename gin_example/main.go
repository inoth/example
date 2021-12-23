package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.New()

	r.Use(func(c *gin.Context) {
		// 自定义异常捕获中间件
		defer func() {
			if err := recover(); err != nil {
				/*
					switch e := err.(type) {
						case BusinessErr:
						// 自定义错误类型
						c.String(500,e.Error())
						default:
						// 错误通用返回
					}
				*/
				c.String(500, err.(error).Error())
				// 中断管道执行，直接返回
				c.Abort()
			}
		}()
		c.Next()
	})

	// 添加一个中间件
	r.Use(func(c *gin.Context) {
		// do something
		c.Next()
	})

	group := r.Group("/groupa" /*, 给组单独添加中间件*/)
	{
		group.GET("", func(c *gin.Context) {
			c.String(200, "ok")
		})
	}

	r.GET("/", HelloWorld)

	// 指定监听端口，默认8080
	r.Run(":8080")
}

func HelloWorld(c *gin.Context) {
	c.String(200, "hello world")
}
