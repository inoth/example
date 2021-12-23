package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserState struct {
	Uid  string
	Name string
}

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
		ctx := context.WithValue(c.Request.Context(), "RequestId", randomId())
		c.Set("RequestId", ctx)
		c.Set("UserState", UserState{
			Uid:  "123123",
			Name: "inoth",
		})
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
	// 虚假的链路跟踪方式，使用context参数
	ctx := GetRequestId(c)
	fmt.Printf("外部 %v\n", ctx.Value("RequestId"))
	go func(context.Context) {
		fmt.Printf("groutene 内部 %v\n", ctx.Value("RequestId"))
	}(ctx)

	// 获取登录状态中的用户信息
	userState := GetUserState(c)
	fmt.Printf("%v:%v\n", userState.Uid, userState.Name)

	c.String(200, "hello world")
}

func randomId() string {
	id, _ := uuid.NewUUID()
	return id.String()[:8]
}

func GetRequestId(c *gin.Context) context.Context {
	ctx, _ := c.Get("RequestId")
	return ctx.(context.Context)
}

func GetUserState(c *gin.Context) UserState {
	userState, _ := c.Get("UserState")
	return userState.(UserState)
}
