package main

import (
	"container/heap"
	"fmt"
	"net"
	"sync"
)

type server struct {
	ip       net.IP
	port     uint16
	protocol byte
}

func (k server) Equal(other server) bool {
	return k.ip.Equal(other.ip) && k.port == other.port && k.protocol == other.protocol
}

const (
	_vless = iota
	_vmess
	_trojan
	_ss
	_ssr
)

type Port struct {
	value    uint16
	priority uint16
	index    int
}
type PortPriorityQueue []*Port

func (pq PortPriorityQueue) Len() int { return len(pq) }

func (pq PortPriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PortPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PortPriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Port)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PortPriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // для избежания утечки памяти
	item.index = -1 // для безопасности
	*pq = old[0 : n-1]
	return item
}

func (pq *PortPriorityQueue) update(item *Port, value uint16, priority uint16) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

// MinHeap - структура для управления минимальной кучей
type MinHeap struct {
	pq PortPriorityQueue
	mu sync.Mutex
}

// NewMinHeap - создание новой минимальной кучи
func NewMinHeap() *MinHeap {
	pq := make(PortPriorityQueue, 0)
	heap.Init(&pq)
	return &MinHeap{pq: pq}
}

// Insert - вставка нового элемента в кучу
func (mh *MinHeap) Insert(value uint16) {
	mh.mu.Lock()
	defer mh.mu.Unlock()
	item := &Port{
		value:    value,
		priority: value,
	}
	heap.Push(&mh.pq, item)
}

// ExtractMin - извлечение минимального элемента из кучи
func (mh *MinHeap) ExtractMin() (uint16, error) {
	mh.mu.Lock()
	defer mh.mu.Unlock()
	if mh.pq.Len() == 0 {
		return 0, fmt.Errorf("heap is empty")
	}
	item := heap.Pop(&mh.pq).(*Port)
	return item.value, nil
}

// PeekMin - получение минимального элемента без его удаления
func (mh *MinHeap) PeekMin() (uint16, error) {
	mh.mu.Lock()
	defer mh.mu.Unlock()
	if mh.pq.Len() == 0 {
		return 0, fmt.Errorf("heap is empty")
	}
	return mh.pq[0].value, nil
}

type infoPackageRow struct {
	Id   uint16
	Name string
	Url  string
	Use  bool
}
