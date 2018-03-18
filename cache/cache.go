// Package cache is a short-lived in-memory cache.
package cache

import (
	"sync"
	"time"
)

// A Value is some cached value.
type Value interface{}

// A key is a unique key for a cached value.
type key string

func (c key) String() string {
	return string(c)
}

// An item is one cached value and its metadata.
type item struct {
	value   Value
	expires time.Time
}

// expired returns true if the cache item is expired.
func (c *item) expired() bool {
	return time.Now().After(c.expires)
}

// A Bucket contains cached items.
type Bucket struct {
	ttl   time.Duration
	items map[key]*item

	mu   sync.RWMutex
	quit chan struct{}
	tag  func(k key, t time.Time)
}

// New creates a new cache bucket with the configured time-to-live.
func New(ttl time.Duration) *Bucket {
	b := &Bucket{
		ttl:   ttl,
		items: make(map[key]*item),
		mu:    sync.RWMutex{},
		quit:  make(chan struct{}),
	}
	exp := newExpunger(b)
	go exp.processTags()
	go exp.expungeExpiredEntries()
	b.tag = exp.tag
	return b
}

// Shutdown signals the cache to clean up and quit.
func (c *Bucket) Shutdown() error {
	close(c.quit)
	return nil
}

// Get returns the value stored for the given key or nil and false.
func (c *Bucket) Get(ky string) (Value, bool) {
	k := key(ky)
	c.mu.RLock()
	it, ok := c.items[k]
	c.mu.RUnlock()
	if !ok {
		return nil, false

	} else if it.expired() {
		c.remove(k)
		return nil, false
	}
	return it.value, true
}

// Put writes the given value to the given key.
func (c *Bucket) Put(ky string, val Value) {
	c.mu.Lock()
	defer c.mu.Unlock()
	expiry := time.Now().Add(c.ttl)
	k := key(ky)
	c.items[k] = &item{
		value:   val,
		expires: expiry,
	}
	c.tag(k, expiry)
}

// Memoize uses the cache to store the result of vfn to avoid repeating
// relatively expensive operations for short periods.
func (c *Bucket) Memoize(ky string, vfn func() (Value, error)) (v Value, err error) {
	var ok bool
	if v, ok = c.Get(ky); ok {
		return v, nil
	}
	defer func() {
		if err == nil {
			c.Put(ky, v)
		}
	}()
	return vfn()
}

// remove removes all the keys from the cache.
func (c *Bucket) remove(keys ...key) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, k := range keys {
		delete(c.items, k)
	}
}

// expunger tracks items in the cache and removes expired values from
// the cache at regular intervals.
type expunger struct {
	b *Bucket

	interval time.Duration

	recs map[time.Time][]key
	tags chan tag
	mu   sync.Mutex
}

// tag represents an expiring key.
type tag struct {
	K key
	T time.Time
}

const expungeInterval = 60 * time.Second

func newExpunger(b *Bucket) *expunger {
	return &expunger{
		b:        b,
		interval: expungeInterval,
		recs:     make(map[time.Time][]key),
		tags:     make(chan tag, 10),
		mu:       sync.Mutex{},
	}
}

func (c *expunger) tag(k key, t time.Time) {
	// Use a channel to avoid blocking tags.
	c.tags <- tag{k, t}
}

func (c *expunger) processTags() {
	for {
		select {
		case t := <-c.tags:
			tick := t.T.Truncate(c.interval).Add(c.interval)
			c.mu.Lock()
			c.recs[tick] = append(c.recs[tick], t.K)
			c.mu.Unlock()

		case <-c.b.quit:
			return
		}
	}
}

func (c *expunger) expungeExpiredEntries() {
	for {
		tick := time.Now().Truncate(c.interval)
		c.mu.Lock()
		for t, ks := range c.recs {
			if tick.After(t) {
				c.b.remove(ks...)
				delete(c.recs, t)
			}
		}
		c.mu.Unlock()
		select {
		case <-time.After(c.interval):
			continue

		case <-c.b.quit:
			return
		}
	}
}
