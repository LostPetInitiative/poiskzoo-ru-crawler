package utils

import (
	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/types"
)

// An IntHeap is a max-heap of CardIDs.
type CardIDHeap []types.CardID

func (h CardIDHeap) Len() int           { return len(h) }
func (h CardIDHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h CardIDHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *CardIDHeap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(types.CardID))
}

func (h *CardIDHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
