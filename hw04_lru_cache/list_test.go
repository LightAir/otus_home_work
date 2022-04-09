package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())

		otherList := NewList()
		l.Remove(otherList.PushFront(1))

		require.Equal(t, 0, l.Len())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})

	t.Run("move to front one", func(t *testing.T) {
		l := NewList()
		node := l.PushBack(5)
		l.MoveToFront(node)

		require.Equal(t, 5, l.Front().Value)
	})

	t.Run("move to front two", func(t *testing.T) {
		l := NewList()
		l.PushBack(5)
		node := l.PushBack(10)

		l.MoveToFront(node)

		require.Equal(t, 10, l.Front().Value)
	})

	t.Run("push string", func(t *testing.T) {
		l := NewList()
		l.PushFront("Test string")

		require.Equal(t, "Test string", l.Front().Value)
	})

	t.Run("push front test", func(t *testing.T) {
		l := NewList()

		firstPush := l.PushFront(10) // [10]
		l.PushFront(20)              // [20, 10]
		l.PushFront(15)              // [15, 20, 10]
		l.PushFront(2)               // [2, 15, 20, 10]
		lastPush := l.PushFront(500) // [500, 2, 15, 20, 10]

		require.Equal(t, 5, l.Len())
		require.Equal(t, l.Front(), lastPush)
		require.Equal(t, l.Back(), firstPush)

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{500, 2, 15, 20, 10}, elems)
	})

	t.Run("push back test", func(t *testing.T) {
		l := NewList()

		firstPush := l.PushBack(10) // [10]
		l.PushBack(20)              // [10, 20]
		l.PushBack(15)              // [10, 20, 15]
		l.PushBack(2)               // [10, 20, 15, 2]
		lastPush := l.PushBack(500) // [10, 20, 15, 2, 500]

		require.Equal(t, 5, l.Len())
		require.Equal(t, l.Back(), lastPush)
		require.Equal(t, l.Front(), firstPush)

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{10, 20, 15, 2, 500}, elems)
	})

	t.Run("test one element", func(t *testing.T) {
		l := NewList()
		lItem := l.PushFront(17)

		require.Equal(t, 17, l.Back().Value)

		l.Remove(lItem)
		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Back())
		require.Nil(t, l.Front())
	})
}
