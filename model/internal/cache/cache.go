package cache

import (
	"sync"
	"time"
)

type Value interface{}

type key string

func (c key) String() string {
	return string(c)
}

type item struct {
	value   Value
	expires time.Time
}

func (c *item) expired() bool {
	return time.Now().After(c.expires)
}

type Bucket struct {
	items map[key]*item
	mu    *sync.RWMutex
	ttl   time.Duration
	quit  chan struct{}
	tag   func(k key, t time.Time)
}

func New(ttl time.Duration) *Bucket {
	b := &Bucket{
		items: make(map[key]*item),
		mu:    &sync.RWMutex{},
		ttl:   ttl,
		quit:  make(chan struct{}),
	}
	exp := newExpunger(b)
	go exp.processTags()
	go exp.expungeExpiredEntries()
	b.tag = exp.tag
	return b
}

func (c *Bucket) Shutdown() error {
	close(c.quit)
	return nil
}

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

func (c *Bucket) remove(keys ...key) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, k := range keys {
		delete(c.items, k)
	}
}

type expunger struct {
	b *Bucket

	interval time.Duration

	recs map[time.Time][]key
	tags chan tag
	mu   *sync.Mutex
}

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
		mu:       &sync.Mutex{},
	}
}

func (c *expunger) tag(k key, t time.Time) {
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
		for t, ks := range c.recs {
			if tick.After(t) {
				c.b.remove(ks...)
				c.mu.Lock()
				delete(c.recs, t)
				c.mu.Unlock()
			}
		}
		select {
		case <-time.After(c.interval):
			continue

		case <-c.b.quit:
			return
		}
	}
}
