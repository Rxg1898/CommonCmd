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
