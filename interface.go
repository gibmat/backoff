package backoff

import (
	"context"
	"math/rand"
	"sync"
	"time"
)

const (
	optkeyFactor       = "factor"
	optkeyInterval     = "interval"
	optkeyJitterFactor = "jitter-factor"
	optkeyMaxRetries   = "max-retries"
	optkeyThreshold    = "threshold"
)

type CancelFunc func()

type Policy interface {
	Start(context.Context) (Backoff, CancelFunc)
}

type Backoff interface {
	Done() <-chan struct{}
	Next() <-chan struct{}
}

type Constant struct {
	delay      time.Duration
	maxRetries int
}

type Option interface {
	Name() string
	Value() interface{}
}

type baseBackoff struct {
	callCount  int
	cancelFunc context.CancelFunc
	ctx        context.Context
	maxRetries int
	mu         sync.Mutex
	next       chan struct{}
}

type constantBackoff struct {
	*baseBackoff
	policy *Constant
}

// Exponential implements an exponential backoff policy.
type Exponential struct {
	factor       float64
	interval     time.Duration
	jitterFactor float64
	maxRetries   int
	random       *rand.Rand
	threshold    time.Duration // max backoff
}

type exponentialBackoff struct {
	*baseBackoff
	policy  *Exponential
	attempt float64
}
