package local_cache

import "errors"

var (
	ErrKeyNotExist = errors.New("该key不存在")
	ErrKeyExpired = errors.New("该key已过期")
	ErrKeyValue = errors.New("该key断言失败")
)
