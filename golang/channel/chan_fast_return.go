package main

/*
需求目的:
	一些场景下, 只需要最快的结果, 慢的结果不关心 , 比查询接口时 , 只需要接口最快返回即可
	实现获取最快返回结果的需求
*/
import (
	"fmt"
	"time"

	"github.com/relunctance/goutils/fc"
)

func main() {
	start := time.Now()
	n := 10
	h := newHandler()
	ret := h.worker(n)
	// 统计耗时
	fmt.Printf("ret:[%d] cost time: %s\n", ret, time.Now().Sub(start).String())
}

// 用于控制返回和实现
type handler struct {
	stopCh chan struct{}
	reqCh  chan int
}

func newHandler() *handler {
	h := &handler{}
	h.stopCh = make(chan struct{})
	h.reqCh = make(chan int)
	return h
}

// 收到停止后, 不在处理请求
func (h *handler) Stop() {
	close(h.stopCh)
}

func (h *handler) worker(n int) int {
	for i := 0; i < n; i++ {
		go func(id int) {
			val := doWoker(id)
			for {
				select {
				case h.reqCh <- val:
					return
				// 当接收到停止的时候,直接停止
				case <-h.stopCh:
					return
				}
			}
		}(i)
	}

	for {
		select {
		case val := <-h.reqCh:
			h.Stop() // 只要第一个有值, 直接返回
			return val
		}
	}

}

// 实际工作函数
func doWoker(id int) int {
	// randNum := fc.Rand(1, 10)
	randNum := fc.RandGenerator(10)
	time.Sleep(time.Duration(randNum) * time.Second)
	fmt.Printf("id:[%d] , randNum:[%d]\n", id, randNum)
	return id * 100
}
