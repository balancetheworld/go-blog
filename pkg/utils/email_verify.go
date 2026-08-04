package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"strings"
	"time"
	"crypto/subtle"
)

// 每条邮箱验证码只存两个信息：验证码、过期时间。
type emailVerifyValue struct {
	Code      string
	ExpiresAt time.Time
}

// emailVerifyStore 全局共享的验证码内存存储容器（匿名结构体变量）
// 内嵌 sync.Mutex 互斥锁：提供 Lock() / Unlock()，防止并发读写map冲突
var emailVerifyStore = struct {
	sync.Mutex                             // 匿名结构体嵌入互斥锁，自带加解锁方法
	values     map[string]emailVerifyValue // key：用户邮箱字符串；value：验证码+过期时间
}{
	values: make(map[string]emailVerifyValue), // 初始化空map，必须手动make分配内存
}

func RequestEmailVerify(email string) string {
	value, err := rand.Int(
		rand.Reader,
		big.NewInt(1000000),
	)
	if err != nil {
		return ""
	}

	code := fmt.Sprintf("%06d", value.Int64())
	key := normalizeEmail(email)

	emailVerifyStore.Lock()
	emailVerifyStore.values[key] = emailVerifyValue{
		Code:      code,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	emailVerifyStore.Unlock()

	return code
}

func normalizeEmail(email string) string {
        return strings.ToLower(
                strings.TrimSpace(email),
        )
  }

  func VerifyEmailCode(email, code string) bool {
	key := normalizeEmail(email)

	emailVerifyStore.Lock()
	defer emailVerifyStore.Unlock()
	 value, exists := emailVerifyStore.values[key]
        if !exists {
                return false
        }

        if time.Now().After(value.ExpiresAt) {
                delete(emailVerifyStore.values, key)
                return false
        }

        if subtle.ConstantTimeCompare(
                []byte(value.Code),
                []byte(code),
        ) != 1 {
                return false
        }

        delete(emailVerifyStore.values, key)
        return true
  }