package cache_test

import (
	"testing"
	"time"

	"github.com/motki/core/model/internal/cache"
)

func Test(t *testing.T) {
	expected := "value"

	b := cache.New(10 * time.Millisecond)
	defer func() {
		if err := b.Shutdown(); err != nil {
			t.Errorf("error shutting down bucket: %s", err.Error())
		}
	}()
	b.Put("test", 0, expected)

	v, ok := b.Get("test", 0)
	if !ok {
		t.Errorf("expected value from cache, got nothing")
		return
	}
	actual, ok := v.(string)
	if !ok {
		t.Errorf("expected value to be string, got %T", v)
		return
	}
	if actual != expected {
		t.Errorf("expected \"%s\", got \"%s\"", expected, actual)
	}

	<-time.After(20 * time.Millisecond)

	_, ok = b.Get("test", 0)
	if ok {
		t.Errorf("expected value to be expunged from cache")
		return
	}
}

func TestMemoize(t *testing.T) {
	b := cache.New(10 * time.Millisecond)
	defer func() {
		if err := b.Shutdown(); err != nil {
			t.Errorf("error shutting down bucket: %s", err.Error())
		}
	}()

	calls := new(int)

	get := func() (cache.Value, error) {
		return b.Memoize("test", 0, func() (cache.Value, error) {
			*calls++
			return *calls, nil
		})
	}

	for i := 0; i < 5; i++ {
		v, err := get()
		if err != nil {
			t.Errorf("error getting value from cache: %s", err.Error())
			return
		}
		vi, ok := v.(int)
		if !ok {
			t.Errorf("expected int, got %T", v)
			return
		}
		if vi != 1 {
			t.Errorf("expected func to be called once, but was called %d times", vi)
			return
		}
	}

	if *calls != 1 {
		t.Errorf("expected func to be called once, but was called %d times", *calls)
		return
	}

	<-time.After(20 * time.Millisecond)

	v, err := get()
	if err != nil {
		t.Errorf("error getting value from cache: %s", err.Error())
		return
	}
	vi, ok := v.(int)
	if !ok {
		t.Errorf("expected int, got %T", v)
		return
	}
	if vi != 2 {
		t.Errorf("expected func to be called twice, but was called %d times", vi)
		return
	}
}
