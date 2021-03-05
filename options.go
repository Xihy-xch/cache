package local_cache

import "time"

type Options struct {
	expiration time.Duration
	mode       Mode
	maxSum     int64
}

type Mode int8

const (
	Default Mode = iota
	LRU
)

type OptionsFn func(opts *Options)

func WithExpiration(expiration time.Duration) OptionsFn {
	return func(opts *Options) {
		opts.expiration = expiration
	}
}

func WithMode(mode Mode) OptionsFn {
	return func(opts *Options) {
		opts.mode = mode
	}
}

func WithMaxSum(maxSum int64) OptionsFn {
	return func(opts *Options) {
		opts.maxSum = maxSum
	}
}

func (o *Options) GetExpiration() time.Duration {
	if o != nil {
		return o.expiration
	}
	return 0
}
