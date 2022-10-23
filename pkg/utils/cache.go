package utils

type kvp[TKey comparable, TData any] struct {
	key   TKey
	value TData
}

type LRUCache[TKey comparable, TData any] struct {
	list     *LinkedList[kvp[TKey, TData]]
	Capacity int
	dict     map[TKey]*ListNode[kvp[TKey, TData]]
}

func NewLRUCache[TKey comparable, TData any](capacity int) *LRUCache[TKey, TData] {
	return &LRUCache[TKey, TData]{
		list:     NewLinkedList[kvp[TKey, TData]](),
		Capacity: capacity,
		dict:     make(map[TKey]*ListNode[kvp[TKey, TData]]),
	}
}

func (c *LRUCache[TKey, TData]) Set(key TKey, val TData) {
	old, exists := c.dict[key]
	if exists {
		c.list.Remove(old)
	}
	node := NewLinkedListNode(kvp[TKey, TData]{key, val})
	c.list.PushAsFirst(node)
	c.dict[key] = node
	for c.list.Size() > c.Capacity {
		toRemove := c.list.Last
		c.list.RemoveLast()
		delete(c.dict, toRemove.Data.key)
	}
}

// bool - whether the extraction is successful, thus first value of tuple is properly set
func (c *LRUCache[TKey, TData]) Get(key TKey) (*TData, bool) {
	node, exists := c.dict[key]
	if exists {
		// moving to head
		c.list.Remove(node)
		c.list.PushAsFirst(node)
		return &node.Data.value, true
	}
	return nil, false

}
