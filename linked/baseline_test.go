package linked

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNew(t *testing.T) {
	n := NewLinkedList(10)
	require.NotNil(t, n)
}

func TestValues(t *testing.T) {
	n := NewLinkedList(5)
	v := n.Values()
	require.Equal(t, []int{1, 2, 3, 4, 5}, v)
}

func TestAppend(t *testing.T) {
	n := NewLinkedList(5)
	require.Equal(t, 0, n.start)
	require.Equal(t, 0, n.count)
	for i := 0; i < 3; i++ {
		err := n.Append(20)
		require.NoError(t, err)
	}
	require.Equal(t, 0, n.start)
	require.Equal(t, 3, n.count)
	require.Equal(t, []int{20, 20, 20, 4, 5}, n.Values())
	for i := 0; i < 2; i++ {
		err := n.Append(20)
		require.NoError(t, err)
	}
	require.Equal(t, 0, n.start)
	require.Equal(t, 5, n.count)
	err := n.Append(120)
	require.Error(t, err)

	require.Equal(t, []int{20, 20, 20, 20, 20}, n.Values())

}

func TestDeleteCount(t *testing.T) {
	n := NewLinkedList(5)
	require.Equal(t, []int{1, 2, 3, 4, 5}, n.Values())
	require.Equal(t, 0, n.start)
	require.Equal(t, 0, n.count)
	n.DeleteCount(3)
	require.Equal(t, 3, n.start)
	require.Equal(t, []int{0, 0, 0, 4, 5}, n.Values())
	require.Equal(t, 0, n.count)
	n.DeleteCount(5)
	require.Equal(t, []int{0, 0, 0, 0, 0}, n.Values())
	require.Equal(t, 3, n.start)
	require.Equal(t, 0, n.count)
}

func TestDeleteBounds(t *testing.T) {
	n := NewLinkedList(5)
	require.Equal(t, []int{1, 2, 3, 4, 5}, n.Values())
	require.Equal(t, 0, n.start)
	require.Equal(t, 0, n.count)
	n.DeleteBounds(0, 0) // delete single value
	require.Equal(t, 1, n.start)
	require.Equal(t, 0, n.count)
	require.Equal(t, []int{0, 2, 3, 4, 5}, n.Values())
	n.DeleteBounds(0, 4)
	require.Equal(t, []int{0, 0, 0, 0, 0}, n.Values())
	require.Equal(t, 0, n.start)
	require.Equal(t, 0, n.count)
	n = NewLinkedList(5)
	require.Equal(t, 0, n.start)
	require.Equal(t, 0, n.count)
	n.DeleteBounds(4, 1)
	require.Equal(t, 2, n.start)
	require.Equal(t, 0, n.count)
}

func TestGetIndex(t *testing.T) {
	n := NewLinkedList(10)
	first := n.GetIndex(0)
	last := n.GetIndex(9)
	fmt.Println(first, last)
	require.Equal(t, 1, first.Value.(int))
	require.Equal(t, 10, last.Value.(int))
}

func TestExample(t *testing.T) {
	n := NewLinkedList(5)
	newStart := 10
	for i := 0; i < 5; i++ {
		err := n.Append(newStart + i)
		require.NoError(t, err)
	}

	n.DeleteBounds(3, 0)
	require.Equal(t, []int{0, 11, 12, 0, 0}, n.Values())
	require.Equal(t, 1, n.start)
	n.Append(100)
	require.Equal(t, []int{0, 11, 12, 100, 0}, n.Values())
	for i := 0; i < 2; i++ {
		err := n.Append(200 + 1 + i)
		require.NoError(t, err)
	}
	err := n.Append(-1)
	require.Error(t, err)

	require.Equal(t, []int{202, 11, 12, 100, 201}, n.Values())
	require.Equal(t, 1, n.start)
	require.Equal(t, 5, n.count)

}
