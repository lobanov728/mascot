package command

import (
	"context"
	"errors"
	"time"
)

var ErrWrongInterfaceType = errors.New("wrong interface type")

type Lock interface {
	Release(ctx context.Context) error
}

type Locker interface {
	AcuireLock(ctx context.Context, key string, ttl time.Duration) (Lock, error)
	Release(ctx context.Context, lock Lock) error
}
