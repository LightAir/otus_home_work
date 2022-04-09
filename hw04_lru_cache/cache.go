package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	lock     sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.lock.Lock()
	defer l.lock.Unlock()

	lItem := l.items[key]

	if lItem == nil {
		return nil, false
	}

	cItem := lItem.Value.(cacheItem)
	l.queue.MoveToFront(lItem)

	return cItem.value, true
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.lock.Lock()
	defer l.lock.Unlock()

	isKeyExist := l.items[key] != nil

	if isKeyExist {
		l.queue.Remove(l.items[key])
	} else if l.queue.Len() >= l.capacity {
		cItem := l.queue.Back().Value.(cacheItem)

		l.queue.Remove(l.items[cItem.key])
		l.items[cItem.key] = nil
	}

	l.items[key] = l.queue.PushFront(cacheItem{key: key, value: value})

	return isKeyExist
}

func (l *lruCache) Clear() {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}
