package lru

import "fmt"

type Cacher interface {
	Get(interface{}) (interface{}, error)
	Put(interface{}, interface{}) error
}

type lruCache struct {
	size      int               // size of cache
	remaining int               // remaining capacity
	cache     map[string]string // actual storage of data
	queue     []string          // for keeping track of lru
}

func NewCache(size int) Cacher {
	return &lruCache{size: size, remaining: size, cache: make(map[string]string), queue: make([]string, 0)}
}

func (lru *lruCache) Get(key interface{}) (interface{}, error) {
	// get(key) -> value, error

	k := (key.(string))

	if len(lru.queue) == lru.size { // delete head of queue if queue == size of cache
		lru.qDel(lru.queue[0])
	}

	lru.qDel(k)                      // delete any copies of k
	lru.queue = append(lru.queue, k) // append k to tail of queue

	if value, ok := lru.cache[k]; ok {
		return value, nil
	}
	return 0, fmt.Errorf("value not there") // return error if value is not found

}

func (lru *lruCache) Put(key, val interface{}) error {
	// put(key, value)

	v := (val.(string)) // key and val are both strings
	k := (key.(string))

	if lru.remaining <= 0 { // delete lru of queue and val of cache associated with queue[0]
		delete(lru.cache, lru.queue[0])
		lru.qDel(lru.queue[0])
		lru.remaining++
	}

	if lru.remaining > 0 { // check for space in cache
		lru.qDel(k) // delete copies of k in queue
		lru.queue = append(lru.queue, k)
		lru.cache[k] = v
		lru.remaining-- // decrease remaining size

	}

	fmt.Println("Queue : ", lru.queue)
	fmt.Println("Cache : ", lru.cache)
	fmt.Println("Remaining Space: ", lru.remaining)

	return nil
}

// Delete element from queue
func (lru *lruCache) qDel(ele string) {
	for i := 0; i < len(lru.queue); i++ {
		if lru.queue[i] == ele {
			oldlen := len(lru.queue)
			copy(lru.queue[i:], lru.queue[i+1:])
			lru.queue = lru.queue[:oldlen-1]
			break
		}
	}
}
