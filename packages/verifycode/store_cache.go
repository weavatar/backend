package verifycode

import (
	"time"

	"github.com/goravel/framework/facades"
)

// CacheStore 实现 verifycode.Store interface
type CacheStore struct {
	KeyPrefix string
}

// Set 实现 verifycode.Store interface 的 Set 方法
func (s *CacheStore) Set(key string, value string) bool {

	ExpireTime := time.Minute * time.Duration(int64(facades.Config().GetInt("verifycode.expire_time")))
	// 本地环境方便调试
	if facades.Config().GetBool("app.debug") {
		ExpireTime = time.Minute * time.Duration(int64(facades.Config().GetInt("verifycode.debug_expire_time")))
	}

	err := facades.Cache().Put(s.KeyPrefix+key, value, ExpireTime)

	return err == nil
}

// Get 实现 verifycode.Store interface 的 Get 方法
func (s *CacheStore) Get(key string, clear bool) (value string) {
	key = s.KeyPrefix + key
	val := facades.Cache().Get(key, "")
	if clear {
		facades.Cache().Forget(key)
	}
	return val.(string)
}

// Verify 实现 verifycode.Store interface 的 Verify 方法
func (s *CacheStore) Verify(key, answer string, clear bool) bool {
	v := s.Get(key, clear)
	return v == answer
}
