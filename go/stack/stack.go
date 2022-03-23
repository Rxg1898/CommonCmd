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
