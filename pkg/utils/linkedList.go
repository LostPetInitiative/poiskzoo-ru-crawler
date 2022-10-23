package utils

type ListNode[TData any] struct {
	Next *ListNode[TData]
	Prev *ListNode[TData]
	Data TData
}

type LinkedList[TData any] struct {
	First *ListNode[TData]
	Last  *ListNode[TData]
	size  int
}

// Constructs an empty list
func NewLinkedList[TData any]() *LinkedList[TData] {
	return &LinkedList[TData]{}
}

func (l *LinkedList[TData]) Size() int {
	return l.size
}

func NewLinkedListNode[TData any](data TData) *ListNode[TData] {
	return &ListNode[TData]{
		Data: data,
	}
}

// pushes the node to the head of the list. The entire list will be attached to the new head
func (l *LinkedList[TData]) PushAsFirst(newFirst *ListNode[TData]) {
	if newFirst == nil {
		panic("can't use nil node")
	}
	prevFirst := l.First
	l.First = newFirst

	newFirst.Next = prevFirst
	newFirst.Prev = nil

	if prevFirst != nil {
		prevFirst.Prev = newFirst
	}

	if l.Last == nil {
		l.Last = newFirst
	}

	l.size++
}

// cuts out the specified element from the list
func (l *LinkedList[TData]) Remove(nodeToRemove *ListNode[TData]) {
	prev := nodeToRemove.Prev
	next := nodeToRemove.Next

	if nodeToRemove == l.First {
		l.First = next
	}
	if nodeToRemove == l.Last {
		l.Last = prev
	}

	if prev != nil {
		prev.Next = next
	}
	if next != nil {
		next.Prev = prev
	}

	l.size--
}

func (l *LinkedList[TData]) RemoveLast() {
	if l.Size() == 0 {
		panic("Can't remove the last node of the list as the list is empty")
	}

	l.Remove(l.Last)
}
