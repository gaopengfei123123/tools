package topKFrequent

// 非常 nb 的写法, 收藏学习一下

type Sortable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64 | ~string
}

type HeapSortOrder int

const (
	ASC HeapSortOrder = iota
	DESC
)

type Heap[T Sortable] struct {
	elems []T
	order HeapSortOrder
}

func NewHeap[T Sortable](order HeapSortOrder) *Heap[T] {
	h := &Heap[T]{
		elems: make([]T, 0),
		order: order,
	}
	return h
}

func (h *Heap[T]) IsEmpty() bool {
	return h.Size() == 0
}

func (h *Heap[T]) Size() int {
	return len(h.elems)
}

func (h *Heap[T]) Peek() (T, bool) {
	if h.IsEmpty() {
		var t T
		return t, false
	}

	return h.elems[0], true
}

func (h *Heap[T]) Push(item T) bool {
	h.elems = append(h.elems, item)
	h.shiftUp(len(h.elems) - 1)
	return true
}

func (h *Heap[T]) Pop() (T, bool) {
	if h.IsEmpty() {
		var t T
		return t, false
	}

	elem := h.elems[0]
	h.swap(0, h.Size()-1)
	h.elems = h.elems[:h.Size()-1]
	h.shiftDown(0)
	return elem, true
}

func (h *Heap[T]) Elems() []T {
	return h.elems
}

func (h *Heap[T]) leftLeafIndex(i int) int {
	return 2*i + 1
}

func (h *Heap[T]) rightLeafIndex(i int) int {
	return 2*i + 2
}

func (h *Heap[T]) parentIndex(i int) int {
	return (i - 1) / 2
}

func (h *Heap[T]) swap(i, j int) {
	h.elems[i], h.elems[j] = h.elems[j], h.elems[i]
}

func (h *Heap[T]) compare(a, b T) bool {
	if h.order == ASC {
		return a < b
	} else {
		return a > b
	}
}

func (h *Heap[T]) shiftUp(i int) {
	if i == 0 {
		return
	}

	for {
		j := h.parentIndex(i)
		if j < 0 || !h.compare(h.elems[i], h.elems[j]) {
			break
		}

		h.swap(i, j)
		i = j
	}
}

func (h *Heap[T]) shiftDown(i int) {
	for {
		l := h.leftLeafIndex(i)
		r := h.rightLeafIndex(i)
		dest := i
		if l < h.Size() && h.compare(h.elems[l], h.elems[dest]) {
			dest = l
		}

		if r < h.Size() && h.compare(h.elems[r], h.elems[dest]) {
			dest = r
		}

		if dest == i {
			break
		}

		h.swap(i, dest)
		i = dest
	}
}

func combine(freq, num int) int {
	if num >= 0 {
		return freq<<16 | num
	}

	return freq<<16 | 0x8000 | -num
}

func split(elem int) (freq int, num int) {
	freq = elem >> 16
	num = elem & 0xffff

	if num >= 0x8000 {
		num &= 0x7fff
		num = -num
	}

	return freq, num
}

func topKFrequent(nums []int, k int) []int {
	frequency := make(map[int]int, len(nums)/2)
	for _, num := range nums {
		frequency[num]++
	}

	h := NewHeap[int](ASC)
	for num, freq := range frequency {
		if h.Size() == k {
			elem, _ := h.Peek()
			f, _ := split(elem)

			if freq > f {
				h.Pop()
				h.Push(combine(freq, num))
			}

			continue
		}

		h.Push(combine(freq, num))
	}

	res := make([]int, 0, k)
	for i := 0; i < k; i++ {
		elem, _ := h.Pop()
		_, num := split(elem)
		res = append(res, num)
	}

	return res
}
