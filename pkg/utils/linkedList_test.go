package utils_test

import (
	"testing"

	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/utils"
)

func TestOneElemList(t *testing.T) {
	l := utils.NewLinkedList[int]()

	l.PushAsFirst(utils.NewLinkedListNode(20))

	if l.First.Data != 20 {
		t.Error("First elem is invalid")
		t.Fail()
	}

	if l.Last.Data != 20 {
		t.Error("Last elem is invalid")
		t.Fail()
	}

	if l.Size() != 1 {
		t.Error("Size is invalid")
		t.Fail()
	}
}

func TestPushes(t *testing.T) {
	l := utils.NewLinkedList[int]()

	l.PushAsFirst(utils.NewLinkedListNode(20))
	l.PushAsFirst(utils.NewLinkedListNode(30))
	l.PushAsFirst(utils.NewLinkedListNode(-1))

	if l.First.Data != -1 {
		t.Error("First elem is invalid")
		t.Fail()
	}

	if l.Last.Data != 20 {
		t.Error("Last elem is invalid")
		t.Fail()
	}

	if l.Size() != 3 {
		t.Error("Size is invalid")
		t.Fail()
	}
}

func TestRemoveLast(t *testing.T) {
	l := utils.NewLinkedList[int]()

	l.PushAsFirst(utils.NewLinkedListNode(20))
	l.PushAsFirst(utils.NewLinkedListNode(30))
	l.PushAsFirst(utils.NewLinkedListNode(-1))
	l.RemoveLast()

	if l.First.Data != -1 {
		t.Error("First elem is invalid")
		t.Fail()
	}

	if l.Last.Data != 30 {
		t.Error("Last elem is invalid")
		t.Fail()
	}

	if l.Size() != 2 {
		t.Error("Size is invalid")
		t.Fail()
	}
}

func TestRemoveFistAsSpecific(t *testing.T) {
	l := utils.NewLinkedList[int]()

	l.PushAsFirst(utils.NewLinkedListNode(20))
	l.PushAsFirst(utils.NewLinkedListNode(30))
	n := utils.NewLinkedListNode(-1)
	l.PushAsFirst(n)
	l.Remove(n)

	if l.First.Data != 30 {
		t.Error("First elem is invalid")
		t.Fail()
	}

	if l.Last.Data != 20 {
		t.Error("Last elem is invalid")
		t.Fail()
	}

	if l.Size() != 2 {
		t.Error("Size is invalid")
		t.Fail()
	}
}

func TestRemoveLastAsSpecific(t *testing.T) {
	l := utils.NewLinkedList[int]()

	n := utils.NewLinkedListNode(-1)
	l.PushAsFirst(n)
	l.PushAsFirst(utils.NewLinkedListNode(20))
	l.PushAsFirst(utils.NewLinkedListNode(30))
	l.Remove(n)

	if l.First.Data != 30 {
		t.Error("First elem is invalid")
		t.Fail()
	}

	if l.Last.Data != 20 {
		t.Error("Last elem is invalid")
		t.Fail()
	}

	if l.Size() != 2 {
		t.Error("Size is invalid")
		t.Fail()
	}
}

func TestMiddleRemove(t *testing.T) {
	l := utils.NewLinkedList[int]()

	l.PushAsFirst(utils.NewLinkedListNode(20))
	middle := utils.NewLinkedListNode(30)
	l.PushAsFirst(middle)
	l.PushAsFirst(utils.NewLinkedListNode(-1))

	l.Remove(middle)

	if l.First.Data != -1 {
		t.Error("First elem is invalid")
		t.Fail()
	}

	if l.Last.Data != 20 {
		t.Error("Last elem is invalid")
		t.Fail()
	}

	if l.Size() != 2 {
		t.Error("Size is invalid")
		t.Fail()
	}
}

func TestRemoveAllButMiddle(t *testing.T) {
	l := utils.NewLinkedList[int]()

	l.PushAsFirst(utils.NewLinkedListNode(20))
	l.PushAsFirst(utils.NewLinkedListNode(30))
	n := utils.NewLinkedListNode(-1)
	l.PushAsFirst(n)
	l.Remove(n)
	l.RemoveLast()

	if l.First.Data != 30 {
		t.Error("First elem is invalid")
		t.Fail()
	}

	if l.Last.Data != 30 {
		t.Error("Last elem is invalid")
		t.Fail()
	}

	if l.Size() != 1 {
		t.Error("Size is invalid")
		t.Fail()
	}
}

func TestRemoveAll(t *testing.T) {
	l := utils.NewLinkedList[int]()

	l.PushAsFirst(utils.NewLinkedListNode(20))
	l.PushAsFirst(utils.NewLinkedListNode(30))
	n := utils.NewLinkedListNode(-1)
	l.PushAsFirst(n)
	l.RemoveLast()
	l.RemoveLast()
	l.RemoveLast()

	if l.First != nil {
		t.Error("First elem is invalid")
		t.Fail()
	}

	if l.Last != nil {
		t.Error("Last elem is invalid")
		t.Fail()
	}

	if l.Size() != 0 {
		t.Error("Size is invalid")
		t.Fail()
	}
}
