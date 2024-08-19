package cache

import "sync"

type Node struct {
	key   string
	value interface{}
	prev  *Node
	next  *Node
}

type Cache struct {
	capacity int
	ll       *Node
	tail     *Node
	cache    map[string]*Node
	mutex    sync.RWMutex
}

func NewCache(capacity int) *Cache {
	return &Cache{
		capacity: capacity,
		ll:       nil,
		tail:     nil,
		cache:    make(map[string]*Node),
		mutex:    sync.RWMutex{},
	}
}

func (c *Cache) moveToHead(node *Node) {
	if c.ll == node {
		return
	}

	if c.tail == node {
		c.tail = node.prev
		c.tail.next = nil
	} else {
		node.prev.next = node.next
		node.next.prev = node.prev
	}
	node.prev = nil
	node.next = c.ll
	c.ll.prev = node
	c.ll = node
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	node, ok := c.cache[key]
	if !ok {
		return nil, false
	}
	c.moveToHead(node)

	return node.value, true
}

func (c *Cache) Set(key string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if node, ok := c.cache[key]; ok {
		node.value = value
		c.moveToHead(node)
		return
	}

	newNode := &Node{
		key:   key,
		value: value,
	}

	c.cache[key] = newNode

	if c.ll == nil {
		c.ll = newNode
		c.tail = newNode
	} else {
		newNode.next = c.ll
		c.ll.prev = newNode
		c.ll = newNode
	}

	if len(c.cache) > c.capacity {
		c.Evict()
	}
}

func (c *Cache) Evict() {
	if c.tail == nil {
		return
	}

	delNode := c.tail
	c.tail = delNode.prev
	if c.tail != nil {
		c.tail.next = nil
	}

	delete(c.cache, delNode.key)
}
