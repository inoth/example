package main

import (
	"fmt"
	"sync"
	"time"
)

// telegf 采集器数据传输简化模型

type Msg struct {
	Value string
}

type in struct {
	dst chan<- Msg
}

type mid struct {
	src <-chan Msg
	dst chan<- Msg
}

type out struct {
	src <-chan Msg
}

func main() {
	next, out := output()

	next, md := middleware(next)
	// fmt.Printf("mid 输出: %p\n", next)
	next, md1 := middleware(next)

	in := input(next)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		runOut(out)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		runMid(md)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		runMid(md1)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		runIn(in)
	}()

	wg.Wait()
	fmt.Printf("结束模拟流程\n")
}

// 初始化输出
func output() (chan<- Msg, *out) {
	src := make(chan Msg, 100)
	fmt.Printf("out 输出: %p\n", src)
	return src, &out{src: src}
}

// 初始化中中间处理
// 给出一组中间处理程序
func middleware(dst chan<- Msg) (chan<- Msg, *mid) {
	src := make(chan Msg, 100)
	fmt.Printf("mid 输出: %p\n", src)
	// 正常状态下, 会存在多份处理器, 每次处理完成后需要把结果对 dst 进行一次覆盖, 避免下一个处理器处理时进度重置
	md := &mid{
		src: src,
		dst: dst,
	}

	// dst = src
	return src, md
}

// 初始化输入
// 给定一个输入列表
func input(dst chan<- Msg) *in {
	return &in{
		dst: dst,
	}
}

func runOut(ou *out) {
	for val := range ou.src {
		fmt.Printf("输出内容: %v\n", val.Value)
	}
}

func runMid(m *mid) {
	for val := range m.src {
		val.Value = val.Value + " tags"
		m.dst <- val
	}
}

func runIn(i *in) {
	// i.dst <- Msg{Value: "test"}
	for {
		i.dst <- Msg{Value: time.Now().String()}
		time.Sleep(time.Second * 1)
	}
}
