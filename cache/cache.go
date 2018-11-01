package cache

import "github.com/zendesk/goship/resources"

// GoshipCache is interface that every cache should implement
type GoshipCache interface {
	CacheName() string
	Resources() []resources.Resource
	Refresh(bool) error
	Save() error
	Read() error
	Len() int
}
