# 数据结构之栈

生活中查看微信聊天信息 ，一般都是从最新时间点往后看(`最新的消息, 最先处理`)

像查看微信聊天信息这种数据结构：FILO(First In Last Out)，先进后出。这种结构称为栈

栈stack，它是一种运算受限的线性表。限定仅在表尾进行插入和删除操作的线性表，这一端被称为栈顶, 另一端是封死的.



## 设计栈？

栈2个核心方法

- Push
- Pop

要满足先进后出这个条件, 这种有序的元素存储我们可以选择数组或者slice, 因此slice支持更多操作，我们选择slice作为存储数据的容器

定义栈这种数据结构

```go
package stack

// 定义需要存入的元素类型
// 这里Item是泛型，指代任意类型
type Item interface{}

// 构造函数
func NewStack() *Stack {
	return &Stack{
		items: []Item{},
	}
}

type Stack struct {
	items []Item
}

// 实现Push,把slice append的那边作为栈顶, 就可以直接放进去
func (s *Stack) Push(item Item) {
	s.items = append(s.items, item)
}

// 实现Pop,我们把栈顶的元素弹出来, 对于slice来说 就是删除最后一个元素, 并把删除的那个元素返回
func (s *Stack) Pop() Item {
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item
}
```

测试用例

```go
package stack_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go/stack"
)

func TestStack(t *testing.T) {
	should := assert.New(t)

	s := stack.NewStack()
	s.Push(1)

	// 判断返回 和 预期是否相等
	should.Equal(1, s.Pop())
}
```



## 完善栈数据结构

上面的程序有那些问题:

- 如果栈到底了 使用Pop会出现指针 (边界问题处理)
- 改数据结构不是线程安全 (并发资源竞争问题)
- slice共用底层数组的的问题，Pop并不能真正的删除元素, 其占用的内存并不会减少当Pop时



为了让栈更丰满，需要补充一些辅助方法:

- 栈大小 Len
- 栈的容量 Size
- 判断是否为空 IsEmpty
- 判断是否已满 IsFull
- 获取栈顶元素的值 Peek
- 清空栈 Clear
- 查询某个值==最近==距离栈顶的距离 Search
- 遍历栈 ForEach

```go
package stack

import (
	"fmt"
	"sync"
)

// 定义需要存入的元素类型
// 这里Item是泛型，指代任意类型
type Item interface{}

// 构造函数
func NewStack() *Stack {
	return &Stack{
		items: []Item{},
	}
}

type Stack struct {
	sync.Mutex
	items []Item
}

// 实现Push,把slice append的那边作为栈顶, 就可以直接放进去
func (s *Stack) Push(item Item) {
	s.Lock()
	s.items = append(s.items, item)
	s.Unlock()
}

// 实现Pop,我们把栈顶的元素弹出来, 对于slice来说 就是删除最后一个元素, 并把删除的那个元素返回
func (s *Stack) Pop() Item {

	if s.IsEmpty() {
		return nil
	}

	item := s.items[len(s.items)-1]
	s.Lock()
	s.items = s.items[:len(s.items)-1]
	s.Unlock()
	return item
}

func (s *Stack) Len() int {
	return len(s.items)
}

func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

func (s *Stack) Peek() Item {

	if s.IsEmpty() {
		return nil
	}

	return s.items[len(s.items)-1]
}

func (s *Stack) Clear() {
	s.items = []Item{}
}

func (s *Stack) ForEach(fn func(Item)) {
	for i := range s.items {
		fn(i)
	}
}

func (s *Stack) Search(item Item) (pos int, err error) {
	for i := range s.items {
		if item == s.items[i] {
			return i, nil
		}
	}
	return 0, fmt.Errorf("item %s not found", item)
}
```

测试用例

```go
package stack_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go/stack"
)

func TestStack(t *testing.T) {
	should := assert.New(t)

	s := stack.NewStack()
	s.Push(1)
	s.Push(2)
	s.Push("a")

	// 判断返回 和 预期是否相等
	should.Equal("a", s.Pop())
	should.Equal(2, s.Pop())
	should.Equal(1, s.Pop())
}
```



## 栈实现插入排序

```
1. 我们把值往另一边放, 大的放下面, 小的放上面 (找大数)
2. 如果发现右边的栈顶没左边的大, 就把右边的栈顶值挪去左边, 直到右边栈顶最大
3. 持续这个循环, 直到左边为空(没值了)
4. 将右边的值 倒过来给左边, 就完成了大的在顶部 小的在底部 大->小 排序完成 
```

```go
// 插入排序 栈顶从大到小排序
func (s *Stack) Sort() {
	// 一个辅助stack
	orderStack := NewStack()

	for !s.IsEmpty() {
		// 开始排序
		tmp := s.Pop()

		// 当前元素大于右边, 应该把右边的挪过左边, 直到右边是最大的元素
		for !orderStack.IsEmpty() && tmp.(int) > orderStack.Peek().(int) {
			s.Push(orderStack.Pop())
		}
		// 当前值一定小于等于右边的值
		orderStack.Push(tmp)
	}

	// 倒排回来
	for !orderStack.IsEmpty() {
		s.Push(orderStack.Pop())
	}
}
```

测试用例

```go
func TestStackOrder(t *testing.T) {
	should := assert.New(t)

	s := stack.NewStack()
	s.Push(1)
	s.Push(2)
	s.Push(5)
	s.Push(3)
	
	// 栈顶从大到小排序
	s.Sort()

	//  判断返回 和 预期是否相等
	should.Equal(5, s.Pop())
	should.Equal(3, s.Pop())
	should.Equal(2, s.Pop())
	should.Equal(1, s.Pop())
}
```

## [示例代码当前目录](./stack.go)

