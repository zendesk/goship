package cache

import (
	"fmt"
	"github.com/zendesk/goship/config"
	"sync"
	"time"

	"github.com/zendesk/goship/color"
)

// GlobalCacheList holds list of caches
type GlobalCacheList []GoshipCache

// RefreshInParallel function refreshes all caches from GlobalCacheList
func (cl *GlobalCacheList) RefreshInParallel(force bool) error {
	var wg sync.WaitGroup
	wg.Add(len(*cl))
	for _, c := range *cl {
		go func(cache GoshipCache) {
			refreshStart := time.Now()
			defer wg.Done()
			refreshed, err := cache.Refresh(force)
			if err != nil {
				color.PrintRed(fmt.Sprintf("Error while refreshing cache: %s\n", err))
			}
			refreshElapsed := time.Since(refreshStart)
			if config.GlobalConfig.Verbose {
				if refreshed {
					color.PrintGreen(fmt.Sprintf("[%s] Cache refresh took: %s. Read %d instances\n", cache.CacheName(), refreshElapsed, cache.Len()))
				} else {
					color.PrintGreen(fmt.Sprintf("[%s] Cache read took: %s. Read %d instances\n", cache.CacheName(), refreshElapsed, cache.Len()))
				}
			}
		}(c)
	}
	wg.Wait()
	return nil
}
