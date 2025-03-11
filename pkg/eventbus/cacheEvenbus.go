package event

const (
	CacheEventAddCache    = "cache.add"
	CacheEventDeleteCache = "cache.delete"
	CacheEventUpdateCache = "cache.update"
)

type CacheEvent struct {
	Type  string
	Short string
	Long  string
}

type CacheEventBus struct {
	bus chan CacheEvent
}

func NewCacheEventBus() *CacheEventBus {
	return &CacheEventBus{
		bus: make(chan CacheEvent),
	}
}

func (e *CacheEventBus) Publis(event CacheEvent) {
	e.bus <- event
}

func (e *CacheEventBus) Subscribe() <-chan CacheEvent {
	return e.bus
}
