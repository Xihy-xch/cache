package local_cache

import "time"

type Options struct {
	expiration time.Duration
	cleanMode  CleanMode
	maxSum     int64
}

type CleanMode int8

const (
	Default CleanMode = iota
	LRU
)

type OptionsFn func(opts *Options)

func WithTimeout(expiration time.Duration) OptionsFn {
	return func(opts *Options) {
		opts.expiration = expiration
	}
}

func WithCleanMode(cleanMode CleanMode) OptionsFn {
	return func(opts *Options) {
		opts.cleanMode = cleanMode
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
