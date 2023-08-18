package captcha

import (
	"time"

	"github.com/goravel/framework/facades"
)

type CacheStore struct {
	KeyPrefix string
}

func (s *CacheStore) Set(key string, value string) error {
	ExpireTime := time.Minute * time.Duration(facades.Config().GetInt("captcha.expire_time"))

	if facades.Config().GetBool("app.debug") {
		ExpireTime = time.Minute * time.Duration(facades.Config().GetInt("captcha.debug_expire_time"))
	}

	err := facades.Cache().Put(s.KeyPrefix+key, value, ExpireTime)
	if err != nil {
		return err
	}

	return nil
}

func (s *CacheStore) Get(key string, clear bool) string {
	key = s.KeyPrefix + key
	val := facades.Cache().Get(key, "").(string)
	if clear {
		facades.Cache().Forget(key)
	}
	return val
}

func (s *CacheStore) Verify(key, answer string, clear bool) bool {
	v := s.Get(key, clear)
	return v == answer
}
