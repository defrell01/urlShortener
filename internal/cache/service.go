package cache

import (
	"log"
	"time"
	"urlshortener/pkg/db"
	event "urlshortener/pkg/eventbus"
)

type CacheServiceDeps struct {
	EventBus    *event.CacheEventBus
	RedisClient *db.Redis
}

type CacheService struct {
	EventBus    *event.CacheEventBus
	RedisClient *db.Redis
}

func NewCacheService(deps *CacheServiceDeps) *CacheService {
	return &CacheService{
		EventBus:    deps.EventBus,
		RedisClient: deps.RedisClient,
	}
}

func (s *CacheService) UpdateCache() {
	for msg := range s.EventBus.Subscribe() {

		if msg.Type == event.CacheEventAddCache {
			short := msg.Short
			if len(short) == 0 {
				log.Fatalln("Bad CacheEventAddCache SHORT data")
			}

			long := msg.Long
			if len(long) == 0 {
				log.Fatalln("Bad CacheEventAddCache LONG data")
			}

			err := s.RedisClient.SetCache(short, long, time.Second*db.DefaultTTL)
			if err != nil {
				log.Fatalf("Err setting cache short: %s, long: %s\n", short, long)
			}
		} else if msg.Type == event.CacheEventDeleteCache {
			short := msg.Short
			if len(short) == 0 {
				log.Fatalln("Bad CacheEventAddCache SHORT data")
			}

			err := s.RedisClient.DeleteCache(short)
			if err != nil {
				log.Fatalf("Err deleting cache short: %s\n", short)
			}
		} else if msg.Type == event.CacheEventUpdateCache {
			short := msg.Short
			if len(short) == 0 {
				log.Fatalln("Bad CacheEventAddCache SHORT data")
			}

			long := msg.Long
			if len(long) == 0 {
				log.Fatalln("Bad CacheEventAddCache LONG data")
			}

			err := s.RedisClient.SetCache(short, long, time.Second*db.DefaultTTL)
			if err != nil {
				log.Fatalf("Err updating cache short: %s, long: %s\n", short, long)
			}
		}

	}
}
