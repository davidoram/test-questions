package queue

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestThreadSafeQ(t *testing.T) {
	q := NewTheadSafeQueue(3)
	require.NotNil(t, q)
	require.Equal(t, 0, q.Len())

	require.True(t, q.Push("a"))
	require.True(t, q.Push("b"))
	require.True(t, q.Push("c"))
	require.False(t, q.Push("d"))

	require.Equal(t, "c", q.Pop())
	require.Equal(t, "b", q.Pop())
	require.Equal(t, "a", q.Pop())
	require.Nil(t, q.Pop())
}

func TestThreadSafeQ_multi_thread(t *testing.T) {
	q := NewTheadSafeQueue(1000)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for range 499 {
			q.Push(1)
		}
	}()
	go func() {
		defer wg.Done()
		for range 501 {
			q.Push(2)
		}
	}()
	wg.Wait()
	require.Equal(t, 1000, q.Len())
	cnt := map[int]int{1: 0, 2: 0}
	for range 1000 {
		val, ok := q.Pop().(int)
		require.True(t, ok)
		cnt[val] = cnt[val] + 1
	}
	require.Equal(t, 2, len(cnt))
	require.Equal(t, 499, cnt[1])
	require.Equal(t, 501, cnt[2])
}
