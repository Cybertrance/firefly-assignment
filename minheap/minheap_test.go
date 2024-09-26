package minheap

import (
	"container/heap"
	"firefly-assignment/utils"
	"testing"
)

func TestMinHeap(t *testing.T) {
	tests := []struct {
		name         string
		initialHeap  []utils.WordFreq
		pushElement  utils.WordFreq
		expectedHeap []utils.WordFreq
		popOrder     []utils.WordFreq
	}{
		{
			name: "Push and Pop single element",
			initialHeap: []utils.WordFreq{
				{Word: "banana", Frequency: 10},
			},
			pushElement: utils.WordFreq{Word: "apple", Frequency: 5},
			expectedHeap: []utils.WordFreq{
				{Word: "apple", Frequency: 5},
				{Word: "banana", Frequency: 10},
			},
			popOrder: []utils.WordFreq{
				{Word: "apple", Frequency: 5},
				{Word: "banana", Frequency: 10},
			},
		},
		{
			name: "Push element to a larger heap",
			initialHeap: []utils.WordFreq{
				{Word: "banana", Frequency: 10},
				{Word: "cherry", Frequency: 20},
			},
			pushElement: utils.WordFreq{Word: "apple", Frequency: 5},
			expectedHeap: []utils.WordFreq{
				{Word: "apple", Frequency: 5},
				{Word: "cherry", Frequency: 20},
				{Word: "banana", Frequency: 10},
			},
			popOrder: []utils.WordFreq{
				{Word: "apple", Frequency: 5},
				{Word: "banana", Frequency: 10},
				{Word: "cherry", Frequency: 20},
			},
		},
		{
			name: "Pop elements from heap",
			initialHeap: []utils.WordFreq{
				{Word: "banana", Frequency: 10},
				{Word: "apple", Frequency: 5},
				{Word: "cherry", Frequency: 20},
			},
			pushElement: utils.WordFreq{Word: "orange", Frequency: 15},
			expectedHeap: []utils.WordFreq{
				{Word: "apple", Frequency: 5},
				{Word: "banana", Frequency: 10},
				{Word: "cherry", Frequency: 20},
				{Word: "orange", Frequency: 15},
			},
			popOrder: []utils.WordFreq{
				{Word: "apple", Frequency: 5},
				{Word: "banana", Frequency: 10},
				{Word: "orange", Frequency: 15},
				{Word: "cherry", Frequency: 20},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := MinHeap(tt.initialHeap)
			heap.Init(&h)

			// Push the new element
			heap.Push(&h, tt.pushElement)

			// Check the resulting heap
			for i, expected := range tt.expectedHeap {
				if h[i] != expected {
					t.Errorf("expected heap[%d] = %v, but got %v", i, expected, h[i])
				}
			}

			// Pop elements and check the order
			for i, expected := range tt.popOrder {
				popped := heap.Pop(&h).(utils.WordFreq)
				if popped != expected {
					t.Errorf("expected pop[%d] = %v, but got %v", i, expected, popped)
				}
			}
		})
	}
}

func TestMinHeapLen(t *testing.T) {
	h := MinHeap{
		{Word: "apple", Frequency: 5},
		{Word: "banana", Frequency: 10},
	}
	expectedLen := 2
	if h.Len() != expectedLen {
		t.Errorf("expected length %d, got %d", expectedLen, h.Len())
	}
}

func TestMinHeapLess(t *testing.T) {
	h := MinHeap{
		{Word: "apple", Frequency: 5},
		{Word: "banana", Frequency: 10},
	}
	if !h.Less(0, 1) {
		t.Errorf("expected apple (5) < banana (10)")
	}
	if h.Less(1, 0) {
		t.Errorf("expected banana (10) > apple (5)")
	}
}

func TestMinHeapSwap(t *testing.T) {
	h := MinHeap{
		{Word: "apple", Frequency: 5},
		{Word: "banana", Frequency: 10},
	}
	h.Swap(0, 1)
	if h[0].Word != "banana" || h[1].Word != "apple" {
		t.Errorf("expected swap to place banana at 0 and apple at 1, got %v", h)
	}
}
