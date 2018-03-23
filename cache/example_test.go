package cache_test

import (
	"fmt"
	"strconv"
	"time"

	"github.com/motki/core/cache"
)

// This example shows how one might use the cache bucket to avoid repeating
// the same expensive call within a short period.
//
// Memoize is used to create a "lazy" closure that is invoked only when the
// cache does not already contain the result. After the real work is done, the
// resulting value is stored in the bucket. After the bucket's configured
// time-to-live has passed, the value is removed from the cache and the next
// invocation is not cached.
func ExampleBucket_Memoize() {
	// Create a new bucket with a 10 second expiration.
	c := cache.New(10 * time.Second)

	// mfib is a memoized fibonacci sequence implementation.
	var mfib func(n int) int
	mfib = func(n int) int {
		v, err := c.Memoize(strconv.Itoa(n), func() (cache.Value, error) {
			fmt.Println(n)
			// Actual fibonacci sequence implementation.
			if n <= 1 {
				return 1, nil
			}
			return mfib(n-1) + mfib(n-2), nil
		})
		if err != nil {
			return 0
		}
		if i, ok := v.(int); ok {
			return i
		}
		return 0
	}

	val := mfib(5)
	fmt.Println("Result:  ", val)
	fmt.Println("Memoized:", mfib(5))
	// Output:
	// 5
	// 4
	// 3
	// 2
	// 1
	// 0
	// Result:   8
	// Memoized: 8
}
