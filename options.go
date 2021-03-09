package local_cache

import "time"

type Getter interface {
	Get(key string) (interface{}, error)
}

type GetterFunc func(key string) (interface{}, error)

func (f GetterFunc) Get(key string) (interface{}, error) {
	return f(key)
}

type Options struct {
	expiration time.Duration
	getter     Getter
}

type OptionsFn func(opts *Options)

func WithExpiration(expiration time.Duration) OptionsFn {
	return func(opts *Options) {
		opts.expiration = expiration
	}
}

func (o *Options) GetExpiration() time.Duration {
	if o != nil {
		return o.expiration
	}
	return 0
}

func WithGetter(getter Getter) OptionsFn {
	return func(opt *Options) {
		opt.getter = getter
	}
}

func getDefaultOptions() *Options {
	return &Options{
		expiration: 10 * time.Second,
	}
}
