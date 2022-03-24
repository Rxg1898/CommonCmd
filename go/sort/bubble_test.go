package sort_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go/sort"
)

// 冒泡排序 测试用例
func TestBubbleSort(t *testing.T) {
	should := assert.New(t)
	rows := []int{55, 5, 4, 1, 6, 12, 7}
	target := sort.BubbleSort(rows)
	should.Equal([]int{1, 4, 5, 6, 7, 12, 55}, target)
}
