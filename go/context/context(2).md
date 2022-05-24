# 1、介绍
项目地址：[https://github.com/oklog/run](https://github.com/oklog/run)
prometheus就是使用这种方式管理多goroutine编排

run.Group是一种管理goroutine组件生命周期的通用机制，它在任何需要将`多个goroutines作为一个单元整体进行协调`的情况下都很有用。

# 2、使用
创建一个零值的run.Group，然后向其添加actors。actor被定义为一对函数：一个执行函数，它应该同步运行；一个中断函数，当被调用时，它应该使执行函数返回。最后，调用Run，它同时运行所有的组件，等待直到第一个actor退出，调用中断函数，最后只有在所有组件都返回后才将控制权返回给调用者。这个通用的API允许调用者对几乎所有可运行的任务进行编排，并为组件实现定义良好的生命周期管理。
## 2.1 简单使用

- g.add 第一个参数为run函数，返回error，要求该函数长时间运行，遇到错误再退出 
   - 具体场景有 for +ticker +ctx
   - 长时间运行的 http 和rpc 用error chan 在g.add中做
- 第二个参数为interrupt函数，作用是退出时做一些清理操作
```go
import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/oklog/run"
)

func main() {
	// 编排开始
	var g run.Group
	ctxAll, cancelAll := context.WithCancel(context.Background())
	fmt.Println(ctxAll)
	{
		// 处理信号退出的channel
		term := make(chan os.Signal, 1)

		signal.Notify(term, os.Interrupt, syscall.SIGTERM)
		cancelC := make(chan struct{})
		g.Add(
			func() error {
				select {
				case <-term:
					fmt.Println("Receive SIGTERM, exiting gracefully...")
					cancelAll()
					return nil
				case <-cancelC:
					fmt.Println("other cancel exiting")
					return nil
				}
			},
			func(err error) {
				close(cancelC)
			},
		)
	}
	g.Run()

}
```
## 2.2 进阶使用
主协程退出触发cancelAll()，通知ctxAll的所有协程触发ctxAll.Done()退出
```go
import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oklog/run"
)

func main() {
	// 编排开始
	var g run.Group
	ctxAll, cancelAll := context.WithCancel(context.Background())
	{
		// 处理信号退出的channel
		term := make(chan os.Signal, 1)

		signal.Notify(term, os.Interrupt, syscall.SIGTERM)
		cancelC := make(chan struct{})
		g.Add(
			func() error {
				select {
				case <-term:
					fmt.Println("Receive SIGTERM, exiting gracefully...")
					cancelAll()
					return nil
				case <-cancelC:
					fmt.Println("other cancel exiting")
					return nil
				}
			},
			func(err error) {
				close(cancelC)
			},
		)
	}
	{
		g.Add(func() error {
			for {
				ticker := time.NewTicker(3 * time.Second)
				select {
				case <-ctxAll.Done():
					fmt.Println("打工人01，接收到了cancelAll的退出指令")
					return nil
				case <-ticker.C:
					fmt.Println("我是打工人01")

				}
			}
		}, func(err error) {

		},
		)
	}
	g.Run()

}
```
# 3、源码分析
## 3.1 run.Group 
![](https://blog-1301758797.cos.ap-guangzhou.myqcloud.com/%E6%96%87%E6%A1%A3%E5%9B%BE%E7%89%87/Go/context/group1.png)

- actor的结构体

![](https://blog-1301758797.cos.ap-guangzhou.myqcloud.com/%E6%96%87%E6%A1%A3%E5%9B%BE%E7%89%87/Go/context/group2.png)
## 3.2 Add

- 向actor组`添加(Add)`一个函数，每个actor必须可以被一个`interrupt中断函数`抢占。也就是说，如果`interrupt中断`被调用，`execute应该返回`。而且，即使在execute返回之后，调用interrupt也必须是安全的。
- 第一个返回的actor（函数）会`中断所有正在运行的actors`
- 错误被传递给中断函数，并由Run返回

![](https://blog-1301758797.cos.ap-guangzhou.myqcloud.com/%E6%96%87%E6%A1%A3%E5%9B%BE%E7%89%87/Go/context/add.png)
## 3.3 Run
`概述`：同时运行所有的actors函数，当第一个actor返回时，所有其他的actors都会被打断。只有当所有的actors都退出时，Run才会返回。`运行返回第一个退出的actor所返回的错误`

`细节`：

- 判断是否有actors，就是有没有Add
- len运行每一个actors
- 等待第一个actor停止
- 向所有actors发出interrupt(中断)信号
- 等待所有actors都停止下来，有些快有些慢所以用channel
- 返回第一个actor触发的错误

![](https://blog-1301758797.cos.ap-guangzhou.myqcloud.com/%E6%96%87%E6%A1%A3%E5%9B%BE%E7%89%87/Go/context/run.png)
