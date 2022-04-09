package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	First *ListItem
	Last  *ListItem
	Size  int
}

func NewList() List {
	return new(list)
}

func (l list) Len() int {
	return l.Size
}

func (l list) Front() *ListItem {
	if l.Size == 0 {
		return nil
	}

	return l.First
}

func (l list) Back() *ListItem {
	if l.Size == 0 {
		return nil
	}

	return l.Last
}

func (l *list) PushFront(v interface{}) *ListItem {
	newNode := &ListItem{
		Value: v,
		Prev:  nil,
		Next:  nil,
	}

	l.Size++

	if l.First == nil {
		l.First = newNode
		l.Last = newNode

		return l.First
	}

	node := l.First

	newNode.Next = node
	if node.Prev == nil {
		l.First = newNode
	}

	node.Prev = newNode

	return newNode
}

func (l *list) PushBack(v interface{}) *ListItem {
	if l.Last == nil {
		return l.PushFront(v)
	}

	node := l.Last

	newNode := &ListItem{
		Value: v,
		Prev:  l.Last,
		Next:  nil,
	}

	if node.Next == nil {
		l.Last = newNode
	}

	node.Next = newNode

	l.Size++

	return newNode
}

func (l *list) Remove(i *ListItem) {
	if l.Size == 0 {
		return
	}

	if i == nil {
		return
	}

	if i.Prev == nil {
		l.First = i.Next
	} else {
		i.Prev.Next = i.Next
	}

	if i.Next == nil {
		l.Last = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}

	l.Size--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.Size <= 1 {
		return
	}

	if i != l.First {
		l.PushFront(i.Value)
		l.Remove(i)
	}
}
