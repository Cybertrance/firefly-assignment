// Package minheap implements a minimum heap (priority queue) for word frequencies.
package minheap

import "firefly-assignment/utils"

// This is a priority queue that keeps the smallest frequency at the top. Used for quite. ach insertion or removal from the heap is O(log N), and we process n elements, so the overall complexity is O(n log N).
// This ensures efficient performance even with a large word frequency map.

type MinHeap []utils.WordFreq

// Len returns the length of the heap.
func (h MinHeap) Len() int { return len(h) }

// Less compares the frequency of two elements, used to order the heap.
func (h MinHeap) Less(i, j int) bool { return h[i].Frequency < h[j].Frequency }

// Swap swaps two elements in the heap.
func (h MinHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

// Push pushes a new word frequency to the heap
func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(utils.WordFreq))
}

// Pop removes the smallest frequency from the heap
func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
