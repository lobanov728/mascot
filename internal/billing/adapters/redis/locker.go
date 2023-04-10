package redis

// "context"

// "github.com/bsm/redislock"

// type service struct {
// 	client *redislock.Client
// }

// func (p *service) Lock(ctx context.Context, key string, timeout, ttl, heartbeat time.Duration) (*redislock.Lock, error) {
// 	backoff := redislock.LinearBackoff(heartbeat)
// 	lockCtx, cancel := context.WithDeadline(ctx, time.Now().Add(timeout))

// 	defer cancel()

// 	// Obtain lock with retry + custom deadline
// 	lock, err := p.client.Obtain(lockCtx, key, ttl, &redislock.Options{RetryStrategy: backoff})
// 	if err != nil {
// 		return nil, errors.Wrap(err, "acquire lock")
// 	}

// 	return lock, nil
// }

// func (p *service) Release(ctx context.Context, lock *redislock.Lock) error {
// 	return lock.Release(ctx)
// }
