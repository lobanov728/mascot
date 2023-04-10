package adapters

import (
	"context"
	"sync"
	"time"

	"github.com/Lobanov728/mascot/internal/billing/app/command"
)

type SimpleLocker struct {
	locked sync.Map
}

type Lock struct {
	key string
}

func (l Lock) Release(ctx context.Context) error {
	return nil
}

func NewSimpleLocker() command.Locker {
	return &SimpleLocker{
		locked: sync.Map{},
	}
}

func (s *SimpleLocker) AcuireLock(ctx context.Context, key string, ttl time.Duration) (command.Lock, error) {
	// The loaded result is true if the value was loaded, false if stored.
	lock := Lock{key: key}
	_, loaded := s.locked.LoadOrStore(key, lock)

	if loaded {
		return nil, nil
	}

	go func(ttl time.Duration) {
		<-time.After(ttl)
		s.Release(context.Background(), &lock)
	}(ttl)

	return &lock, nil
}

func (s *SimpleLocker) Release(ctx context.Context, lock command.Lock) error {
	l, ok := lock.(Lock)
	if !ok {
		return command.ErrWrongInterfaceType
	}
	s.locked.Delete(l.key)

	return nil
}
